package inject

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/hidevopsio/hiboot/pkg/starter/data"
	"github.com/hidevopsio/hiboot/pkg/starter"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"path/filepath"
	"os"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type user struct {
	Name string
	App  string
}

type fakeRepository struct {
	data.BaseKVRepository
}

type fakeProperties struct {
	Name     string
	Nickname string
	Username string
	Url      string
}

type fakeConfiguration struct {
	Properties fakeProperties `mapstructure:"fake"`
}

type fakeDataSource struct {
}

type FakeRepository data.Repository

type FooUser struct {
	Name string
}

func (c *fakeConfiguration) FakeRepository() FakeRepository {
	repo := new(fakeRepository)
	repo.SetDataSource(new(fakeDataSource))
	return repo
}

// FooUser an instance fooUser is injectable with tag `inject:"fooUser"`
func (c *fakeConfiguration) FooUser() *FooUser {
	u := new(FooUser)
	u.Name = "foo"
	return u
}

type fooConfiguration struct {
}

type fooService struct {
	FooUser       *FooUser       `inject:"name=foo"` // TODO: should be able to change the instance name, e.g. `inject:"bazUser"`
	FooRepository FakeRepository `inject:""`
}

type hibootService struct {
	HibootUser *user `inject:"name=${app.name}"`
}

type barService struct {
	FooRepository FakeRepository `inject:""`
}

type userService struct {
	FooUser        *FooUser       `inject:"name=foo"`
	User           *user          `inject:""`
	FakeUser       *user          `inject:"name=${fake.name},app=${app.name}"`
	FakeRepository FakeRepository `inject:""`
	DefaultUrl     string         `value:"${fake.defaultUrl:http://localhost:8080}"`
	Url            string         `value:"${fake.url}"`
}

type sliceInjectionTestService struct {
	Profiles         []string       `value:"${app.profiles.include}"`
}

type fooBarService struct {
	FooBarRepository FakeRepository `inject:""`
}

type foobarRecursiveInject struct {
	FoobarService *fooBarService `inject:""`
}

type recursiveInject struct {
	UserService *userService
}

type MethodInjectionService struct {
	fooUser    *FooUser
	barUser    *user
	repository FakeRepository
}

var (
	appName    = "hiboot"
	fakeName   = "fake"
	fooName    = "foo"
	fakeUrl    = "http://fake.com/api/foo"
	defaultUrl = "http://localhost:8080"
)

func init() {
	utils.EnsureWorkDir("../..")

	configPath := filepath.Join(utils.GetWorkDir(), "config")
	fakeFile := "application-fake.yaml"
	os.Remove(filepath.Join(configPath, fakeFile))
	fakeContent :=
		"fake:" +
			"\n  name: " + fakeName +
			"\n  nickname: ${app.name} ${fake.name}\n" +
			"\n  username: ${unknown.name:bar}\n" +
			"\n  url: " + fakeUrl
	utils.WriterFile(configPath, fakeFile, []byte(fakeContent))

	starter.Add("fake", fakeConfiguration{})
	starter.Add("foo", fooConfiguration{})
	starter.GetAutoConfiguration().Build()
}

// Init automatically inject FooUser and FakeRepository that instantiated in fakeConfiguration
func (s *MethodInjectionService) Init(fooUser *FooUser, barUser *user, repository FakeRepository) {
	s.fooUser = fooUser
	s.barUser = barUser
	s.repository = repository
}

func TestNotInject(t *testing.T) {
	baz := new(userService)
	assert.Equal(t, (*user)(nil), baz.User)
}

func TestInject(t *testing.T) {
	t.Run("should inject through method", func(t *testing.T) {
		s := new(MethodInjectionService)
		err := IntoObject(reflect.ValueOf(s))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*FooUser)(nil), s.fooUser)
		assert.NotEqual(t, (*user)(nil), s.barUser)
		assert.NotEqual(t, (FakeRepository)(nil), s.repository)
	})

	t.Run("should inject repository", func(t *testing.T) {
		us := new(userService)
		err := IntoObject(reflect.ValueOf(us))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*user)(nil), us.User)
		assert.Equal(t, fooName, us.FooUser.Name)
		assert.Equal(t, fakeName, us.FakeUser.Name)
		assert.Equal(t, appName, us.FakeUser.App)
		assert.Equal(t, fakeUrl, us.Url)
		assert.Equal(t, defaultUrl, us.DefaultUrl)
		assert.NotEqual(t, (*fakeRepository)(nil), us.FakeRepository)
	})

	t.Run("should not inject unimplemented interface into FooBarRepository", func(t *testing.T) {
		fb := new(foobarRecursiveInject)
		err := IntoObject(reflect.ValueOf(fb))
		assert.Contains(t, err.Error(), "FakeRepository is not implemented")
	})

	t.Run("should not inject unimplemented interface into FooRepository", func(t *testing.T) {
		fs := new(fooService)
		err := IntoObject(reflect.ValueOf(fs))
		assert.Equal(t, "foo", fs.FooUser.Name)
		assert.Contains(t, err.Error(), "FakeRepository is not implemented")
	})

	t.Run("should not inject system property into object", func(t *testing.T) {
		fs := new(hibootService)
		err := IntoObject(reflect.ValueOf(fs))
		assert.Equal(t, nil, err)
		assert.Equal(t, appName, fs.HibootUser.Name)
	})

	t.Run("should not inject unimplemented interface into BarRepository", func(t *testing.T) {
		bs := new(barService)
		err := IntoObject(reflect.ValueOf(bs))
		assert.Contains(t, err.Error(), "FakeRepository is not implemented")
	})

	t.Run("should inject recursively", func(t *testing.T) {
		ps := &recursiveInject{UserService: new(userService)}
		err := IntoObject(reflect.ValueOf(ps))
		assert.Equal(t, nil, err)
		assert.NotEqual(t, (*user)(nil), ps.UserService.User)
		assert.NotEqual(t, (*fakeRepository)(nil), ps.UserService.FakeRepository)
	})

	t.Run("should not inject slice", func(t *testing.T) {
		testSvc := struct {
			Users []FooUser `inject:""`
		}{}
		err := IntoObject(reflect.ValueOf(testSvc))
		assert.Equal(t, "slice injection is not implemented", err.Error())
	})

	t.Run("should inject slice value", func(t *testing.T) {
		testSvc := new(sliceInjectionTestService)
		err := IntoObject(reflect.ValueOf(testSvc))
		assert.Equal(t, nil, err)
		log.Debug(testSvc.Profiles)
	})
}
