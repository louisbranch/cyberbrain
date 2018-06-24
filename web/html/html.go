package html

import (
	"fmt"
	"html/template"
	"io"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"gitlab.com/luizbranco/cyberbrain/web"
)

type HTML struct {
	basepath string
	sync     sync.RWMutex
	cache    map[string]*template.Template
}

func New(basepath, piioDomain, piioAppID string) *HTML {
	fns["piioScript"] = piioScript(piioDomain, piioAppID)
	fns["img"] = img(piioAppID != "")

	return &HTML{
		basepath: basepath,
		cache:    make(map[string]*template.Template),
	}
}

func (h *HTML) Render(w io.Writer, page web.Page) error {
	paths := append([]string{page.Layout}, page.Partials...)

	for i, n := range paths {
		p := []string{h.basepath}
		p = append(p, strings.Split(n+".html", "/")...)
		paths[i] = filepath.Join(p...)
	}

	tpl, err := h.parse(paths...)
	if err != nil {
		return err
	}

	err = tpl.Execute(w, page)
	return err
}

var fns = template.FuncMap{
	"contains": contains,
	"first":    first,
}

func (h *HTML) parse(names ...string) (tpl *template.Template, err error) {
	cp := make([]string, len(names))
	copy(cp, names)
	sort.Strings(cp)
	id := strings.Join(cp, ":")

	h.sync.RLock()
	tpl, ok := h.cache[id]
	h.sync.RUnlock()

	if !ok {
		tpl = template.New(path.Base(names[0])).Funcs(fns)

		tpl, err = tpl.ParseFiles(names...)
		if err != nil {
			return nil, err
		}
		h.sync.Lock()
		h.cache[id] = tpl
		h.sync.Unlock()
	}

	return tpl, nil
}

func first(list []string) string {
	if len(list) == 0 {
		return ""
	}

	return list[0]
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func piioScript(domain, appID string) func() template.HTML {
	tag := `<script type="application/javascript">
  var piioData = {
    appKey: "%s",
    encodeSrc: false,
    domain: "%s"
  }
</script>
<script src="//js.piio.co/%s/piio.min.js"></script>
`

	if domain != "" && appID != "" {
		tag = fmt.Sprintf(tag, appID, domain, appID)
	} else {
		tag = ""
	}

	return func() template.HTML {
		return template.HTML(tag)
	}
}

func img(optimize bool) func(string) template.HTML {
	tag := `<img src="%s" />`

	if optimize {
		tag = `<img data-piio="%s"
src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8+f9vPQAJZAN2rlRQVAAAAABJRU5ErkJggg==" />`
	}

	return func(imgURL string) template.HTML {
		return template.HTML(fmt.Sprintf(tag, imgURL))
	}
}
