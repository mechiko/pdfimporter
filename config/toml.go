package config

var TomlConfig = []byte(`
# This is a TOML document.
ssccprefix = "146024436369"
ssccstartnumber = 1
perpallet = 24
marktemplate = ""
packtemplate = ""
chunksize = 250
party="Б1"

[layouts]
timelayout = "2006-01-02T15:04:05-0700"
timelayoutclear = "2006.01.02 15:04:05"
timelayoutday = "2006.01.02"
timelayoututc = "2006-01-02T15:04:05"

[application]
license = "f7bc886d-bbcd-4ce9-845f-1209d87d406d"

`)

type Configuration struct {
	SsccPrefix      string `mapstructure:"ssccprefix" toml:"ssccprefix" json:"ssccprefix"`
	SsccStartNumber int    `mapstructure:"ssccstartnumber" toml:"ssccstartnumber" json:"ssccstartnumber"`
	PerPallet       int    `mapstructure:"perpallet" toml:"perpallet" json:"perpallet"`
	MarkTemplate    string `mapstructure:"marktemplate" toml:"marktemplate" json:"marktemplate"`
	PackTemplate    string `mapstructure:"packtemplate" toml:"packtemplate" json:"packtemplate"`
	Party           string `mapstructure:"party" toml:"party" json:"party"`
	ChunkSize       int64  `mapstructure:"chunksize" toml:"chunksize" json:"chunksize"`
	// Application AppConfiguration `mapstructure:"application" toml:"application" json:"application"`
	Layouts LayoutConfiguration `mapstructure:"layouts" toml:"layouts" json:"layouts"`
	// описатели БД рефактор
	// Config    DatabaseConfiguration `json:"config"`
	// AlcoHelp3 DatabaseConfiguration `json:"alcohelp3"`
	// TrueZnak  DatabaseConfiguration `json:"trueznak"`
	// SelfDB    DatabaseConfiguration `json:"selfdb"`
	// описание клиента ЧЗ
	// TrueClient TrueClientConfig `json:"trueclient"`
}

type LayoutConfiguration struct {
	TimeLayout      string `mapstructure:"timelayout" toml:"timelayout" json:"timelayout"`
	TimeLayoutClear string `mapstructure:"timelayoutclear" toml:"timelayoutclear" json:"timelayoutclear"`
	TimeLayoutDay   string `mapstructure:"timelayoutday" toml:"timelayoutday" json:"timelayoutday"`
	TimeLayoutUTC   string `mapstructure:"timelayoututc" toml:"timelayoututc" json:"timelayoututc"`
}

type DatabaseConfiguration struct {
	Connection string `mapstructure:"connection" toml:"connection" json:"connection"`
	Driver     string `mapstructure:"driver" toml:"driver" json:"driver"`
	DbName     string `mapstructure:"dbname" toml:"dbname" json:"dbname"`
	File       string `mapstructure:"file" toml:"file" json:"file"`
	User       string `mapstructure:"user" toml:"user" json:"user"`
	Pass       string `mapstructure:"pass" toml:"pass" json:"pass"`
	Host       string `mapstructure:"host" toml:"host" json:"host"`
	Port       string `mapstructure:"port" toml:"port" json:"port"`
}

type AppConfiguration struct {
	// Pwd          string `json:"pwd"`
	// Console      bool   `json:"console"`
	// Disconnected bool   `json:"disconnected"`
	Fsrarid string `mapstructure:"fsrarid" toml:"fsrarid" json:"fsrarid"`
	// DbType       string `json:"dbtype"`
	License string `mapstructure:"license" toml:"license" json:"license"`
	// ScanTimer    int    `json:"scantimer"`
	StartPage string `mapstructure:"startpage" toml:"startpage" json:"startpage"`
}

type TrueClientConfig struct {
	Test        bool   `mapstructure:"test" toml:"test" json:"test"`
	StandGIS    string `mapstructure:"standgis" toml:"standgis" json:"standgis"`
	StandSUZ    string `mapstructure:"standsuz" toml:"standsuz" json:"standsuz"`
	TestGIS     string `mapstructure:"testgis" toml:"testgis" json:"testgis"`
	TestSUZ     string `mapstructure:"testsuz" toml:"testsuz" json:"testsuz"`
	TokenGIS    string `mapstructure:"tokengis" toml:"tokengis" json:"tokengis"`
	TokenSUZ    string `mapstructure:"tokensuz" toml:"tokensuz" json:"tokensuz"`
	AuthTime    string `mapstructure:"authtime" toml:"authtime" json:"authtime"`
	LayoutUTC   string `mapstructure:"layoututc" toml:"layoututc" json:"layoututc"`
	HashKey     string `mapstructure:"hashkey" toml:"hashkey" json:"hashkey"`
	DeviceID    string `mapstructure:"deviceid" toml:"deviceid" json:"deviceid"`
	OmsID       string `mapstructure:"omsid" toml:"omsid" json:"omsid"`
	UseConfigDB bool   `mapstructure:"useconfigdb" toml:"useconfigdb" json:"useconfigdb"`
}
