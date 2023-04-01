package config

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var Requests = viper.New()
var Brang = viper.New()

var (
	HomeDir, _   = os.UserHomeDir()
	BrangPath    = filepath.Join(HomeDir, "brang")
	ConfigPath   = filepath.Join(BrangPath, "config")
	ConfigFile   = filepath.Join(ConfigPath, "config.yaml")
	RequestsFile = filepath.Join(ConfigPath, "requests.yaml")
)

func LoadBrangConfig() error {
	Brang.SetConfigFile(ConfigFile)
	err := Brang.ReadInConfig()
	if err != nil {
		return err
	}
	Brang.SetDefault("deleteTempFileOnClose", true)
	Brang.SetDefault("outWriterFileName", "brangOutput")
	return nil
}

func LoadRequests() error {
	if err := Requests.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// ignore NotFound. Will return error later if key:value not found
		} else {
			// return all other errors
			return fmt.Errorf("err reading requests: %w", err)
		}
	}
	return nil
}

type output struct {
	format string
}

type stdOutput output

type fileOutput struct {
	output
	file string
}

type tempFileOutput struct {
	output
	deleteOnClose bool
}

type bWriter struct {
	Writer *bufio.Writer
	Format string
	Fn     func()
	Err    error
}

type OutWriter interface {
	Init() *bWriter
}

func (o *stdOutput) Init() *bWriter {
	w := bufio.NewWriter(os.Stdout)
	fn := func() {
		w.Flush()
	}
	return &bWriter{w, o.format, fn, nil}

}

func (o *fileOutput) Init() *bWriter {
	f, err := os.Create(o.file)
	if err != nil {
		return &bWriter{Err: err}
	}
	w := bufio.NewWriter(f)
	fn := func() {
		w.Flush()
		f.Close()
	}
	return &bWriter{w, o.format, fn, nil}
}

func (o *tempFileOutput) Init() *bWriter {
	var fn func()
	var w *bufio.Writer
	f, err := os.CreateTemp("", "brangOutput_*"+OutputFileExt())
	if err != nil {
		return &bWriter{Err: err}
	}
	open := OpenEditor(f.Name())
	w = bufio.NewWriter(f)
	fn = func() {
		w.Flush()
		f.Close()
		open.Run()
		if o.deleteOnClose {
			os.Remove(f.Name())
		}
	}
	return &bWriter{w, o.format, fn, nil}
}

func OutputWriter() OutWriter {
	s := Brang.GetString("outWriterFormat")
	ws := Brang.GetString("outWriter")
	switch ws {
	case "stdout", "":
		return &stdOutput{s}
	case "file":
		p, err := OutWriterFilePath()
		if err != nil {
			fmt.Println(err)
			return &stdOutput{s}
		}
		fname := Brang.GetString("outWriterFileName") + OutputFileExt()
		return &fileOutput{output{s}, filepath.Join(p, fname)}
	case "tempFile":
		return &tempFileOutput{output{s}, Brang.GetBool("deleteTempFileOnClose")}
	default:
		fmt.Println("wrong or missing value for outWriter. Looked for: ", ws)
		return nil
	}
}

func OpenEditor(f string) *exec.Cmd {
	ed := Brang.GetString("fileEditor")
	switch runtime.GOOS {
	case "windows":
		if ed == "" {
			ed = "notepad"
		}
	case "darwin":
		if ed != "" {
			ed = "open -a " + ed
		} else {
			ed = "open"
		}
	default:
		if ed == "" {
			ed = "vim"
		}
	}
	o := exec.Command(ed)
	o.Args = append(o.Args, f)
	return o
}

func OutWriterFilePath() (string, error) {
	p := Brang.GetString("outWriterFilePath")
	if p == "" {
		d := HomeDir
		if d == "" {
			return "", fmt.Errorf("err finding home dir")
		}
		if runtime.GOOS == "windows" {
			d = filepath.Join(d, "Desktop")
		}
		p = d
	}
	return p, nil
}

func OutputFileExt() string {
	ext := Brang.GetString("outWriterFileType")
	if ext == "" {
		ext = ".txt"
	} else {
		ext = "." + ext
	}
	return ext
}
