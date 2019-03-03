//+build gormock

package storage

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/xorcare/gormock"
	"github.com/xorcare/gormock/comp"
	"reflect"
	"testing"
)

func TestStorage_FindByID(t *testing.T) {
	type args struct {
		id uint64
	}

	gk := gormock.New(t)
	gk.
		On("First", &Model{ID: 1}, mock.Anything).
		Run(func(args mock.Arguments) {
			model := args[0].(*Model)
			*model = Model{ID: 1}
		}).Once().
		On("First", &Model{ID: 2}, mock.Anything).
		Return(comp.DB{
			Error: errors.New("error true"),
		}).Once()
	defer gk.AssertExpectations(t)

	storage := Storage{gk.DB()}

	tests := []struct {
		name      string
		repo      Storage
		args      args
		wantModel *Model
		wantErr   bool
	}{
		{
			name: "positive",
			repo: storage,
			args: args{
				id: 1,
			},
			wantModel: &Model{ID: 1},
		},
		{
			name: "negative",
			repo: storage,
			args: args{
				id: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotModel, err := tt.repo.FindByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Storage.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotModel, tt.wantModel) {
				t.Errorf("Storage.FindByID() = %v, want %v", gotModel, tt.wantModel)
			}
		})
	}
}

func TestStorage_FindAll(t *testing.T) {
	gk := gormock.New(t)
	gk.
		On("Find", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			model := args[0].(*[]Model)
			*model = []Model{{ID: 1}, {ID: 2}}
		}).Once().
		On("Find", mock.Anything, mock.Anything).
		Return(comp.DB{
			Error: errors.New("error true"),
		}).Once()
	defer gk.AssertExpectations(t)

	storage := Storage{gk.DB()}

	tests := []struct {
		name       string
		repo       Storage
		wantModels []Model
		wantErr    bool
	}{
		{
			name: "positive",
			repo: storage,
			wantModels: []Model{
				{ID: 1}, {ID: 2},
			},
		},
		{
			name:    "negative",
			repo:    storage,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotModels, err := tt.repo.FindAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotModels, tt.wantModels) {
				t.Errorf("Storage.FindAll() = %v, want %v", gotModels, tt.wantModels)
			}
		})
	}
}

func TestStorage_Count(t *testing.T) {
	gk := gormock.New(t)
	gk.
		On("Model", mock.Anything).Once().
		On("Count", mock.Anything).Once().
		Run(func(args mock.Arguments) {
			model := args[0].(*int64)
			*model = 8
		}).
		On("Model", mock.Anything).Once().
		On("Count", mock.Anything).Once().
		Return(comp.DB{
			Error: errors.New("error true"),
		})
	defer gk.AssertExpectations(t)

	storage := Storage{gk.DB()}

	tests := []struct {
		name       string
		repo       Storage
		wantResult int64
		wantErr    bool
	}{
		{
			name:       "positive",
			repo:       storage,
			wantResult: 8,
		},
		{
			name:       "negative",
			repo:       storage,
			wantResult: 0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := tt.repo.Count()
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("Storage.Count() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestStorage_DeleteByID(t *testing.T) {
	type args struct {
		id uint64
	}

	gk := gormock.New(t)
	gk.
		On("Where", "id = ?", []interface{}{uint64(1)}).Once().
		On("Delete", Model{}, mock.Anything).Once().
		On("Where", "id = ?", []interface{}{uint64(2)}).Once().
		On("Delete", Model{}, mock.Anything).Once().
		Return(comp.DB{
			Error: errors.New("error true"),
		})
	defer gk.AssertExpectations(t)

	storage := Storage{gk.DB()}

	tests := []struct {
		name    string
		repo    Storage
		args    args
		wantErr bool
	}{
		{
			name: "positive",
			repo: storage,
			args: args{
				id: 1,
			},
		},
		{
			name:    "negative",
			repo:    storage,
			wantErr: true,
			args: args{
				id: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.repo.DeleteByID(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_DeleteAll(t *testing.T) {
	gk := gormock.New(t)
	gk.
		On("Delete", Model{}, mock.Anything).Once().
		On("Delete", Model{}, mock.Anything).Once().
		Return(comp.DB{
			Error: errors.New("error true"),
		})
	defer gk.AssertExpectations(t)

	storage := Storage{gk.DB()}

	tests := []struct {
		name    string
		repo    Storage
		wantErr bool
	}{
		{
			name: "positive",
			repo: storage,
		},
		{
			name:    "negative",
			repo:    storage,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.repo.DeleteAll(); (err != nil) != tt.wantErr {
				t.Errorf("Storage.DeleteAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Save(t *testing.T) {
	type args struct {
		model *Model
	}

	gk := gormock.New(t)
	gk.
		On("Save", mock.Anything).Once().
		On("Save", mock.Anything).Once().
		Return(comp.DB{
			Error: errors.New("error true"),
		})
	defer gk.AssertExpectations(t)

	storage := Storage{gk.DB()}

	tests := []struct {
		name    string
		repo    Storage
		args    args
		wantErr bool
	}{
		{
			name: "positive",
			repo: storage,
		},
		{
			name:    "negative",
			repo:    storage,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.repo.Save(tt.args.model); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
