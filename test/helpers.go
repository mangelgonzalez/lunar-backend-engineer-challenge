package test

import (
	"context"
	"lunar-backend-engineer-challenge/cmd/di"
	"lunar-backend-engineer-challenge/internal/rockets/module/rocket/domain"
	test "lunar-backend-engineer-challenge/test/arranger"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var common *di.RocketsDI
var ctx context.Context

func changeDirToRoot() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func setUp() {
	ctx, _ = di.Context()
	if common == nil {
		common = di.InitWithEnvFile("../.env", "../.env.test")
	}
	changeDirToRoot()

	setupResources()
}

func setupResources() {
	arrangeWg := &sync.WaitGroup{}
	arrangeWg.Add(3)

	go test.NewRedisArranger(common.Services).Arrange(ctx, arrangeWg)
	go test.NewMysqlArranger(common).Arrange(ctx, arrangeWg)

	arrangeWg.Wait()
}

func givenTheFollowingRocketExist(t *testing.T, rocket *domain.Rocket) {
	assert.NoError(t, common.RocketServices.Repository.Save(rocket))
}
