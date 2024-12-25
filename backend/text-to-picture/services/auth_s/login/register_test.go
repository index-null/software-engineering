// register_test.go
package auth_s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
	"gorm.io/gorm"

	getDB "text-to-picture/config"
	db "text-to-picture/models/init"

	"text-to-picture/models/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockUserRepository 是一个模拟的用户仓库
type MockUserRepository struct{}

// InsertUserInformation 是一个模拟的插入用户信息方法
func (m *MockUserRepository) InsertUserInformation(user *user.UserInformation) error {
	// 模拟插入成功的场景
	return nil
}

// GetUserByName 是一个模拟的获取用户方法
func (m *MockUserRepository) GetUserByName(name string) (*user.UserInformation, error) {
	// 模拟用户不存在的场景
	return nil, gorm.ErrRecordNotFound
}

// SetupRouter 设置 Gin 路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", Register)
	return r
}

// DBConfig 结构体定义
type DBConfig struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

// TestMain 是测试的入口函数
func TestMain(m *testing.M) {
	// 读取测试数据库配置
	yamlFile, err := os.ReadFile(getDB.GetDBConfigPath())
	if err != nil {
		fmt.Printf("Error reading configs.yaml file: %v\n", err)
		os.Exit(1)
	}

	// 解析配置文件
	var dbconfig DBConfig
	err = yaml.Unmarshal(yamlFile, &dbconfig)
	if err != nil {
		fmt.Printf("Error parsing configs.yaml file: %v\n", err)
		os.Exit(1)
	}

	// 设置数据库连接的环境变量
	os.Setenv("DB_USER", dbconfig.DB.User)
	os.Setenv("DB_PASSWORD", dbconfig.DB.Password)
	os.Setenv("DB_HOST", dbconfig.DB.Host)
	os.Setenv("DB_PORT", dbconfig.DB.Port)
	os.Setenv("DB_NAME", dbconfig.DB.Name)

	// 连接数据库
	if err := db.ConnectDatabase(); err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}

	// 运行测试
	code := m.Run()

	// 清理工作
	// 在这里关闭数据库连接等

	os.Exit(code)
}

// TestRegister_Success 测试注册成功的情况
func TestRegister_Success(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个 POST 请求
	body := bytes.NewBuffer([]byte(`{"email": "test2@example.com", "username": "testuser2", "password": "testpassword"}`)) // 每次测试要修改，避免违反唯一性
	req, _ := http.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", "application/json")

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Println("返回的message为：", response["message"])
	assert.Equal(t, "注册成功", response["message"])
}

// TestRegister_BadRequest 测试请求数据格式错误的情况
func TestRegister_BadRequest(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个 POST 请求，故意发送格式错误的数据
	body := bytes.NewBuffer([]byte(`{"email": "test4@example.com", "username": "testuser4", "password": testpassword}`)) //  password 的值没有用""包围起来
	req, _ := http.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", "application/json")

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 检查响应的message
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Println("返回的message为：", response["message"])
	assert.Equal(t, "请求数据格式错误", response["message"])
}

// TestRegister_InternalServerError 测试用户创建失败的情况
func TestRegister_InternalServerError(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个新的路由器
	router := SetupRouter()

	// 创建一个 POST 请求
	body := bytes.NewBuffer([]byte(`{"email": "testuser1@example.com", "username": "testuser1", "password": "test"}`)) // 密码长度小于6导致创建失败
	req, _ := http.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", "application/json")

	// 创建一个响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 检查响应的message
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Println("返回的message为：", response["message"])
	assert.Equal(t, "用户创建失败", response["message"])
	db.DB.Table("userinformation").Where("username = ?", "testuser1").Delete(&user.UserInformation{})
}
