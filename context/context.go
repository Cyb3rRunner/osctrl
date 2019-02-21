package context

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

const (
	// DefaultEnrollPath as default value for enrolling nodes
	DefaultEnrollPath = "enroll"
	// DefaultLogPath as default value for logging data from nodes
	DefaultLogPath = "log"
	// DefaultConfigPath as default value for configuring nodes
	DefaultConfigPath = "config"
	// DefaultQueryReadPath as default value for distributing on-demand queries to nodes
	DefaultQueryReadPath = "read"
	// DefaultQueryWritePath as default value for collecting results from on-demand queries
	DefaultQueryWritePath = "write"
	// DefaultCarverInitPath as default init endpoint for the carver
	DefaultCarverInitPath = "init"
	// DefaultCarverBlockPath as default block endpoint for the carver
	DefaultCarverBlockPath = "block"
	// DefaultContextIcon as default icon to use for contexts
	DefaultContextIcon = "fas fa-wrench"
	// DefaultContextType as default type of context
	DefaultContextType = "osquery"
	// DefaultSecretLength as default length for secrets
	DefaultSecretLength = 64
)

// TLSContext to hold each of the TLS context
type TLSContext struct {
	gorm.Model
	Name          string `gorm:"index"`
	Hostname      string
	Secret        string
	SecretPath    string
	Type          string
	DebugHTTP     bool
	Icon          string
	Configuration string
	Certificate   string
}

// TLSPath to hold all the paths for TLS
type TLSPath struct {
	EnrollPath      string
	LogPath         string
	ConfigPath      string
	QueryReadPath   string
	QueryWritePath  string
	CarverInitPath  string
	CarverBlockPath string
}

// Context keeps all TLS Contexts
type Context struct {
	DB *gorm.DB
}

// CreateContexts to initialize the context struct
func CreateContexts(backend *gorm.DB) *Context {
	var c *Context
	c = &Context{DB: backend}
	return c
}

// Get TLSContext by name
func (context *Context) Get(name string) (TLSContext, error) {
	var ctx TLSContext
	if err := context.DB.Where("name = ?", name).First(&ctx).Error; err != nil {
		return ctx, err
	}
	return ctx, nil
}

// Empty generates an empty TLSContext with default values
func (context *Context) Empty(name, hostname string) TLSContext {
	return TLSContext{
		Name:          name,
		Hostname:      hostname,
		Secret:        generateRandomString(DefaultSecretLength),
		SecretPath:    generateKSUID(),
		Type:          DefaultContextType,
		DebugHTTP:     false,
		Icon:          DefaultContextIcon,
		Configuration: "",
		Certificate:   "",
	}
}

// Create new TLSContext
func (context *Context) Create(ctx TLSContext) error {
	if context.DB.NewRecord(ctx) {
		if err := context.DB.Create(&ctx).Error; err != nil {
			return fmt.Errorf("Create TLSContext %v", err)
		}
	} else {
		return fmt.Errorf("db.NewRecord did not return true")
	}
	return nil
}

// Exists checks if TLSContext exists already
func (context *Context) Exists(name string) bool {
	var results int
	context.DB.Model(&TLSContext{}).Where("name = ?", name).Count(&results)
	return (results > 0)
}

// All gets all TLSContext
func (context *Context) All() ([]TLSContext, error) {
	var ctxs []TLSContext
	if err := context.DB.Find(&ctxs).Error; err != nil {
		return ctxs, err
	}
	return ctxs, nil
}

// Delete TLSContext by name
func (context *Context) Delete(name string) error {
	ctx, err := context.Get(name)
	if err != nil {
		return fmt.Errorf("error getting context %v", err)
	}
	if err := context.DB.Delete(&ctx).Error; err != nil {
		return fmt.Errorf("Delete %v", err)
	}
	return nil
}

// UpdateConfiguration configuration for a context
func (context *Context) UpdateConfiguration(name, configuration string) error {
	ctx, err := context.Get(name)
	if err != nil {
		return fmt.Errorf("error getting context %v", err)
	}
	if err := context.DB.Model(&ctx).Update("configuration", configuration).Error; err != nil {
		return fmt.Errorf("Update %v", err)
	}
	return nil
}