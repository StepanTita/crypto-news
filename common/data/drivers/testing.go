// //go:build integration

package drivers

import (
	"reflect"
	"runtime"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"common"
	commoncfg "common/config"
	"common/data/model"
	"migrator/migrate"
)

const testConfigPath = "./config.test.yaml"

type IntegrationTest func(t *testing.T, log *logrus.Entry, db *sqlx.DB)

type SuiteIntegrationTest struct {
	t             *testing.T
	cfg           commoncfg.Config
	testInstances []IntegrationTest
}

func NewSuite(t *testing.T) *SuiteIntegrationTest {
	return &SuiteIntegrationTest{
		t: t,
	}
}

// SetupSuite Setup db value
func (t *SuiteIntegrationTest) SetupSuite() {
	common.SetupWorkingDirectory()

	t.cfg = commoncfg.NewFromFile(testConfigPath)

	// Migrate DB up
	require.NoError(t.t, migrate.Migrate(t.cfg, migrate.Up))
}

// TearDownSuite Run After All Test Done (migrate down)
func (t *SuiteIntegrationTest) TearDownSuite() {
	sqlDB := t.cfg.DB()
	defer sqlDB.Close()

	// Migrate DB down
	require.NoError(t.t, migrate.Migrate(t.cfg, migrate.Down))
}

// CleanupSuite Run After All Test Done (cleanup tables)
func (t *SuiteIntegrationTest) CleanupSuite() {
	sqlDB := t.cfg.DB()
	defer sqlDB.Close()
	defer t.TearDownSuite()

	tables := []string{model.NEWS, model.USERS}

	for _, table := range tables {

		_, err := sq.Delete(table).RunWith(t.cfg.DB()).Exec()
		require.NoError(t.t, err)
	}
}

func (t *SuiteIntegrationTest) AddTests(tests ...IntegrationTest) {
	for _, test := range tests {
		t.testInstances = append(t.testInstances, test)
	}
}

func (t *SuiteIntegrationTest) TestRunIntegration() {
	t.cfg.Logging().Debug("Starting integration tests run...")
	for _, testInstance := range t.testInstances {
		testInstance(t.t, t.cfg.Logging().WithField("integration-test", GetFunctionName(testInstance)), t.cfg.DB())
	}
	t.cfg.Logging().Debug("Wrapping up integration tests run...")
}

func GetFunctionName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func TestEqualWithoutFields(t *testing.T, exp, act any, omitFields ...string) {
	vExp := reflect.Indirect(reflect.ValueOf(exp))
	vAct := reflect.Indirect(reflect.ValueOf(act))

	if vExp.Kind() != vAct.Kind() {
		require.Failf(t, "not equal kinds", "expected: %v and actual: %v, are of different kinds - %v, %v", exp, act, vExp.Kind(), vAct.Kind())
	}

	if vExp.Kind() != reflect.Struct {
		require.Failf(t, "not structs", "expected: %v and actual: %v, are not of struct kind - %v, %v", exp, act, vExp.Kind(), vAct.Kind())
	}

	tExp := vExp.Type()
	tAct := vAct.Type()

	if tExp.Name() != tAct.Name() {
		require.Failf(t, "not same structs", "expected: %v and actual: %v, are not the same structs - %s, %s", exp, act, tExp.Name(), tAct.Name())
	}
	for i := 0; i < tExp.NumField(); i++ {
		if slices.Contains(omitFields, tExp.Field(i).Name) {
			continue
		}
		if !reflect.DeepEqual(vExp.Field(i).Interface(), vAct.Field(i).Interface()) {
			require.Failf(t, "not equal", "expected: %v and actual: %v, are not equal fields - %s, %s", vExp.Field(i).Interface(), vAct.Field(i).Interface(), tExp.Field(i).Name, tAct.Field(i).Name)
		}
	}
}
