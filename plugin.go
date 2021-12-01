package plugin

import (
	"fmt"
	"log"

	"github.com/gdbu/scribe"
	"github.com/mojura/kiroku"
	"github.com/mojura/mojura"
	"github.com/vroomy/vroomy"
)

var o opts

func init() {
	if err := vroomy.Register("mojura-opts", &o); err != nil {
		log.Fatal(err)
	}
}

type opts struct {
	vroomy.BasePlugin

	out  *scribe.Scribe
	opts *mojura.Opts

	Source kiroku.Source `vroomy:"mojura-source"`
}

// Load will initialize the s3 client
func (o *opts) Load(env map[string]string) (err error) {
	o.out = scribe.New("Mojura Options")
	if o.opts.Dir = env["mojura-sync-mode"]; len(o.opts.Dir) == 0 {
		o.opts.Dir = env["dataDir"]
	}

	switch env["mojura-sync-mode"] {
	case "development":
		o.out.Notification("Development mode enabled, disabling s3 DB syncing")
	case "mirror":
		o.out.Notification("Mirror mode enabled, enabling s3 read-only DB syncing")
		o.opts.IsMirror = true
		o.opts.Source = o.Source
	case "sync":
		o.out.Notification("Sync mode enabled, enabling s3 DB syncing")
		o.opts.Source = o.Source

	default:
		err = fmt.Errorf("invalid mode, <%s> is not supported, available modes are development, mirror, and sync", env["mojura-sync-mode"])
		return
	}

	return
}

// Backend exposes this plugin's data layer to other plugins
func (o *opts) Backend() interface{} {
	return o.opts
}
