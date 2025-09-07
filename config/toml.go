package config

var TomlConfig = []byte(`
# This is a TOML document.
ssccprefix = "146024436369"
ssccstartnumber = 1
perpallet = 24
marktemplate = ""
packtemplate = ""

[layouts]
timelayout = "2006-01-02T15:04:05-0700"
timelayoutclear = "2006.01.02 15:04:05"
timelayoutday = "2006.01.02"
timelayoututc = "2006-01-02T15:04:05"

[application]
license = "f7bc886d-bbcd-4ce9-845f-1209d87d406d"

`)

type Configuration struct {
	SsccPrefix      string `json:"ssccprefix"`
	SsccStartNumber int    `json:"ssccstartnumber"`
	PerPallet       int    `json:"perpallet"`
	MarkTemplate    string `json:"marktemplate"`
	PackTemplate    string `json:"packtemplate"`
	// Application AppConfiguration    `json:"application"`
	Layouts LayoutConfiguration `json:"layouts"`
	// описатели БД рефактор
	// Config    DatabaseConfiguration `json:"config"`
	// AlcoHelp3 DatabaseConfiguration `json:"alcohelp3"`
	// TrueZnak  DatabaseConfiguration `json:"trueznak"`
	// SelfDB    DatabaseConfiguration `json:"selfdb"`
	// описание клиента ЧЗ
	// TrueClient TrueClientConfig `json:"trueclient"`
}

type LayoutConfiguration struct {
	TimeLayout      string `json:"timelayout"`
	TimeLayoutClear string `json:"timelayoutclear"`
	TimeLayoutDay   string `json:"timelayoutday"`
	TimeLayoutUTC   string `json:"timelayoututc"`
}

type DatabaseConfiguration struct {
	Connection string `json:"connection"`
	Driver     string `json:"driver"`
	DbName     string `json:"dbname"`
	File       string `json:"file"`
	User       string `json:"user"`
	Pass       string `json:"pass"`
	Host       string `json:"host"`
	Port       string `json:"port"`
}

type AppConfiguration struct {
	// Pwd          string `json:"pwd"`
	// Console      bool   `json:"console"`
	// Disconnected bool   `json:"disconnected"`
	Fsrarid string `json:"fsrarid"`
	// DbType       string `json:"dbtype"`
	License string `json:"license"`
	// ScanTimer    int    `json:"scantimer"`
	StartPage string `json:"startpage"`
}

type TrueClientConfig struct {
	Test        bool   `json:"test"`
	StandGIS    string `json:"standgis"`
	StandSUZ    string `json:"standsuz"`
	TestGIS     string `json:"testgis"`
	TestSUZ     string `json:"testsuz"`
	TokenGIS    string `json:"tokengis"`
	TokenSUZ    string `json:"tokensuz"`
	AuthTime    string `json:"authtime"`
	LayoutUTC   string `json:"layoututc"`
	HashKey     string `json:"hashkey"`
	DeviceID    string `json:"deviceid"`
	OmsID       string `json:"omsid"`
	UseConfigDB bool   `json:"useconfigdb"`
}
