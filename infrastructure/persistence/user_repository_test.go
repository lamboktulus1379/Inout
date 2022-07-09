package persistence

import (
	"context"
	"inout/domain/model"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestUserRepository_GetById(t *testing.T) {
	DB, err := NewRepositories()
	if err != nil {
		t.Fatalf("Error connection to database. %v", err)
	}
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser model.User
		wantErr  bool
	}{
		{
			name: "Test Success #1",
			fields: fields{
				DB: DB,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantUser: model.User{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UserRepository{
				DB: tt.fields.DB,
			}
			gotUser, err := repo.GetById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("UserRepository.GetById() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
