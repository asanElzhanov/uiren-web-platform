// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package achievements is a generated GoMock package.
package achievements

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgconn "github.com/jackc/pgx/v5/pgconn"
)

// MockachievementRepo is a mock of achievementRepo interface.
type MockachievementRepo struct {
	ctrl     *gomock.Controller
	recorder *MockachievementRepoMockRecorder
}

// MockachievementRepoMockRecorder is the mock recorder for MockachievementRepo.
type MockachievementRepoMockRecorder struct {
	mock *MockachievementRepo
}

// NewMockachievementRepo creates a new mock instance.
func NewMockachievementRepo(ctrl *gomock.Controller) *MockachievementRepo {
	mock := &MockachievementRepo{ctrl: ctrl}
	mock.recorder = &MockachievementRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockachievementRepo) EXPECT() *MockachievementRepoMockRecorder {
	return m.recorder
}

// addLevel mocks base method.
func (m *MockachievementRepo) addLevel(ctx context.Context, dto AddAchievementLevelDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "addLevel", ctx, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// addLevel indicates an expected call of addLevel.
func (mr *MockachievementRepoMockRecorder) addLevel(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "addLevel", reflect.TypeOf((*MockachievementRepo)(nil).addLevel), ctx, dto)
}

// beginTransaction mocks base method.
func (m *MockachievementRepo) beginTransaction(ctx context.Context) (transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "beginTransaction", ctx)
	ret0, _ := ret[0].(transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// beginTransaction indicates an expected call of beginTransaction.
func (mr *MockachievementRepoMockRecorder) beginTransaction(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "beginTransaction", reflect.TypeOf((*MockachievementRepo)(nil).beginTransaction), ctx)
}

// createAchievement mocks base method.
func (m *MockachievementRepo) createAchievement(ctx context.Context, name string) (achievement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createAchievement", ctx, name)
	ret0, _ := ret[0].(achievement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createAchievement indicates an expected call of createAchievement.
func (mr *MockachievementRepoMockRecorder) createAchievement(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createAchievement", reflect.TypeOf((*MockachievementRepo)(nil).createAchievement), ctx, name)
}

// decrementUpperLevels mocks base method.
func (m *MockachievementRepo) decrementUpperLevels(ctx context.Context, tx transaction, dto DeleteAchievementLevelDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "decrementUpperLevels", ctx, tx, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// decrementUpperLevels indicates an expected call of decrementUpperLevels.
func (mr *MockachievementRepoMockRecorder) decrementUpperLevels(ctx, tx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "decrementUpperLevels", reflect.TypeOf((*MockachievementRepo)(nil).decrementUpperLevels), ctx, tx, dto)
}

// deleteAchievement mocks base method.
func (m *MockachievementRepo) deleteAchievement(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "deleteAchievement", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// deleteAchievement indicates an expected call of deleteAchievement.
func (mr *MockachievementRepoMockRecorder) deleteAchievement(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "deleteAchievement", reflect.TypeOf((*MockachievementRepo)(nil).deleteAchievement), ctx, id)
}

// deleteAchievementLevelsByID mocks base method.
func (m *MockachievementRepo) deleteAchievementLevelsByID(ctx context.Context, achID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "deleteAchievementLevelsByID", ctx, achID)
	ret0, _ := ret[0].(error)
	return ret0
}

// deleteAchievementLevelsByID indicates an expected call of deleteAchievementLevelsByID.
func (mr *MockachievementRepoMockRecorder) deleteAchievementLevelsByID(ctx, achID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "deleteAchievementLevelsByID", reflect.TypeOf((*MockachievementRepo)(nil).deleteAchievementLevelsByID), ctx, achID)
}

// deleteLevel mocks base method.
func (m *MockachievementRepo) deleteLevel(ctx context.Context, tx transaction, dto DeleteAchievementLevelDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "deleteLevel", ctx, tx, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// deleteLevel indicates an expected call of deleteLevel.
func (mr *MockachievementRepoMockRecorder) deleteLevel(ctx, tx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "deleteLevel", reflect.TypeOf((*MockachievementRepo)(nil).deleteLevel), ctx, tx, dto)
}

// getAchievement mocks base method.
func (m *MockachievementRepo) getAchievement(ctx context.Context, id int) (achievement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getAchievement", ctx, id)
	ret0, _ := ret[0].(achievement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getAchievement indicates an expected call of getAchievement.
func (mr *MockachievementRepoMockRecorder) getAchievement(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getAchievement", reflect.TypeOf((*MockachievementRepo)(nil).getAchievement), ctx, id)
}

// getAllAchievements mocks base method.
func (m *MockachievementRepo) getAllAchievements(ctx context.Context) ([]achievement, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getAllAchievements", ctx)
	ret0, _ := ret[0].([]achievement)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getAllAchievements indicates an expected call of getAllAchievements.
func (mr *MockachievementRepoMockRecorder) getAllAchievements(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getAllAchievements", reflect.TypeOf((*MockachievementRepo)(nil).getAllAchievements), ctx)
}

// getLastLevelAndTreshold mocks base method.
func (m *MockachievementRepo) getLastLevelAndTreshold(ctx context.Context, achID int) (LevelData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getLastLevelAndTreshold", ctx, achID)
	ret0, _ := ret[0].(LevelData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getLastLevelAndTreshold indicates an expected call of getLastLevelAndTreshold.
func (mr *MockachievementRepoMockRecorder) getLastLevelAndTreshold(ctx, achID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getLastLevelAndTreshold", reflect.TypeOf((*MockachievementRepo)(nil).getLastLevelAndTreshold), ctx, achID)
}

// getLevel mocks base method.
func (m *MockachievementRepo) getLevel(ctx context.Context, achID, level int) (AchievementLevel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getLevel", ctx, achID, level)
	ret0, _ := ret[0].(AchievementLevel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getLevel indicates an expected call of getLevel.
func (mr *MockachievementRepoMockRecorder) getLevel(ctx, achID, level interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getLevel", reflect.TypeOf((*MockachievementRepo)(nil).getLevel), ctx, achID, level)
}

// getLevelsByAchievementID mocks base method.
func (m *MockachievementRepo) getLevelsByAchievementID(ctx context.Context, achID int) ([]AchievementLevel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getLevelsByAchievementID", ctx, achID)
	ret0, _ := ret[0].([]AchievementLevel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getLevelsByAchievementID indicates an expected call of getLevelsByAchievementID.
func (mr *MockachievementRepoMockRecorder) getLevelsByAchievementID(ctx, achID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getLevelsByAchievementID", reflect.TypeOf((*MockachievementRepo)(nil).getLevelsByAchievementID), ctx, achID)
}

// updateAchievement mocks base method.
func (m *MockachievementRepo) updateAchievement(ctx context.Context, dto UpdateAchievementDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "updateAchievement", ctx, dto)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// updateAchievement indicates an expected call of updateAchievement.
func (mr *MockachievementRepoMockRecorder) updateAchievement(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "updateAchievement", reflect.TypeOf((*MockachievementRepo)(nil).updateAchievement), ctx, dto)
}

// Mocktransaction is a mock of transaction interface.
type Mocktransaction struct {
	ctrl     *gomock.Controller
	recorder *MocktransactionMockRecorder
}

// MocktransactionMockRecorder is the mock recorder for Mocktransaction.
type MocktransactionMockRecorder struct {
	mock *Mocktransaction
}

// NewMocktransaction creates a new mock instance.
func NewMocktransaction(ctrl *gomock.Controller) *Mocktransaction {
	mock := &Mocktransaction{ctrl: ctrl}
	mock.recorder = &MocktransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocktransaction) EXPECT() *MocktransactionMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *Mocktransaction) Commit(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MocktransactionMockRecorder) Commit(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*Mocktransaction)(nil).Commit), ctx)
}

// Exec mocks base method.
func (m *Mocktransaction) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sql}
	for _, a := range arguments {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MocktransactionMockRecorder) Exec(ctx, sql interface{}, arguments ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sql}, arguments...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*Mocktransaction)(nil).Exec), varargs...)
}

// Rollback mocks base method.
func (m *Mocktransaction) Rollback(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MocktransactionMockRecorder) Rollback(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*Mocktransaction)(nil).Rollback), ctx)
}
