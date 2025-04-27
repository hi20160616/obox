package data

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hi20160616/obox/configs"
)

type Object struct {
	Title, Body, Folder, FileTitle string
	Data                           []interface{}
	Err                            error
	Info                           string
}

type Objects struct {
	Title string
	Data  []string
	Err   error
	Info  string
}

type HomePage struct {
	Title, Body, Folder, FileTitle string
	Data                           []interface{}
	Atts                           []os.FileInfo
	Objs                           *Objects
	Err                            error
	Info                           string
}

var (
	objectCache   = make(map[string]*Object)
	objectsCache  *Objects
	homePageCache *Object
)

// pattern like `[!foobar]` means a inter-page need to be made as link
var innerObject = regexp.MustCompile(`\[!.+\]`)

// GetObjects 获取对象列表，优先从缓存中获取
func GetObjects() (*Objects, error) {
	if objectsCache != nil {
		return objectsCache, nil
	}
	var err error
	objectsCache, err = ListObjects()
	return objectsCache, err
}

// GetObject 获取单个对象，优先从缓存中获取
func GetObject(title string) (*Object, error) {
	if obj, ok := objectCache[title]; ok {
		return obj, nil
	}
	var err error
	obj, err := NewObject(title)
	if err != nil {
		return nil, err
	}
	obj, err = Load(obj)
	if err != nil {
		return nil, err
	}
	objectCache[title] = obj
	return obj, nil
}

// GetHomePage 获取主页数据，优先从缓存中获取
func GetHomePage() (*Object, error) {
	if homePageCache != nil {
		return homePageCache, nil
	}
	var err error
	homePageCache, err = LoadHomePage()
	return homePageCache, err
}

// 当对象保存或删除时，需要更新缓存
func UpdateObjectCache(o *Object) {
	objectCache[o.Title] = o
}

func DeleteObjectCache(title string) {
	delete(objectCache, title)
}

// 当对象列表发生变化时，需要更新缓存
func UpdateObjectsCache() error {
	var err error
	objectsCache, err = ListObjects()
	return err
}

func InnerLink(body string) string {
	repl := func(pagename string) string {
		pagename = pagename[2 : len(pagename)-1]
		origin := pagename
		pagename = strings.ReplaceAll(pagename, " ", "-")
		return fmt.Sprintf("[%s](/view/%s)", origin, pagename)
	}

	return innerObject.ReplaceAllStringFunc(body, repl)
}

func LoadHomePage() (*Object, error) {
	o, err := NewObject("Home")
	if err != nil {
		return nil, err
	}
	o, err = Load(o)
	if err != nil {
		return nil, err
	}
	// TODO: innerLink invoke and use in render
	o.Body = InnerLink(o.Body)
	// list home attachments
	files, err := listAttachments(o)
	if err != nil {
		return nil, err
	}
	Atts := []os.FileInfo{}
	// for _, file := range files {
	// 	Atts = append(Atts, file)
	// }
	Atts = append(Atts, files...)
	// list objects
	Objs, err := ListObjects()
	if err != nil {
		return nil, err
	}
	o.Data = append(o.Data, Atts, Objs)
	return o, nil
}

func NewObject(title string) (*Object, error) {
	title, err := url.QueryUnescape(title)
	if err != nil {
		return nil, err
	}
	p := &Object{Title: title}
	p.Folder = filepath.Join(configs.Data.DataPath, title)
	p.FileTitle = filepath.Join(p.Folder, "."+title+".md")
	return p, nil
}

// Save write done Body after NewObject() generate the p
func Save(o *Object) error {
	if _, err := os.Stat(o.Folder); err != nil && os.IsNotExist(err) {
		os.MkdirAll(o.Folder, 0755)
	}
	return os.WriteFile(o.FileTitle, []byte(o.Body), 0600)
}

// load read person info after NewObject() generate the p
func Load(o *Object) (*Object, error) {
	body, err := os.ReadFile(o.FileTitle)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return o, err
		}
		return nil, err
	}
	o.Body = string(body)
	return o, nil
}

func ListObjects() (*Objects, error) {
	objs := &Objects{Title: "Objects list"}
	dirs, err := os.ReadDir(configs.Data.DataPath)
	if err != nil {
		return nil, fmt.Errorf("error walking the path %q: %v", configs.Data.DataPath, err)
	}
	for _, dir := range dirs {
		if dir.IsDir() && strings.ToLower(dir.Name()) != "home" {
			objs.Data = append(objs.Data, dir.Name())
		}
	}
	return objs, nil
}

func GetAttachments(o *Object) ([]fs.FileInfo, error) {
	if configs.Data.RecurseDir {
		return walk(o)
	}
	return readDir(o)
}

func ListAttachments(o *Object) (*Object, error) {
	files, err := GetAttachments(o)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		o.Data = append(o.Data, file)
	}
	return o, nil
}

func listAttachments(o *Object) ([]fs.FileInfo, error) {
	if configs.Data.RecurseDir {
		return walk(o)
	}
	return readDir(o)
}

// walk get all files info in o.Folder
func walk(o *Object) ([]fs.FileInfo, error) {
	files := []fs.FileInfo{}
	err := filepath.WalkDir(o.Folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		// if !d.IsDir() && filepath.Ext(path) != ".md" && d.Name()[:1] != "." {
		if d.Name()[:1] != "." {
			if info, err := d.Info(); err != nil {
				return err
			} else {
				files = append(files, info)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", o.Folder, err)
		return nil, err
	}
	return files, nil
}

func readDir(o *Object) ([]fs.FileInfo, error) {
	files, err := os.ReadDir(o.Folder)
	if err != nil {
		return nil, err
	}
	rt := []fs.FileInfo{}
	for _, file := range files {
		fname := file.Name()
		if file.Type().IsRegular() && filepath.Ext(fname) != ".md" && fname[:1] != "." {
			f, err := file.Info()
			if err != nil {
				return nil, err
			}
			rt = append(rt, f)
		}
	}
	return rt, nil
}
