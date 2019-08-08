// progress-gui will launch a gui window to display the gui of
// a backup.
// progress gui will parse messages from dfb and restic on stdin like the
// following:
//
// {"message_type":"dfb","action":"begin","group":"some-group","domain":"some-domain"}
// {"message_type":"summary","files_new":1,"files_changed":2,"files_unmodified":83,"dirs_new":0,"dirs_changed":0,"dirs_unmodified":0,"data_blobs":0,"tree_blobs":0,"data_added":0,"total_files_processed":83,"total_bytes_processed":43535,"total_duration":0.388768151,"snapshot_id":"xxx"}
//
// gui-progress is based on fyne.io, see:
// https://github.com/fyne-io/fyne

package main

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/nattvara/dfb/internal/gui/progress"
	"github.com/nattvara/dfb/internal/restic"

	"fyne.io/fyne/app"
)

func main() {

	app := app.New()

	p := progress.NewProgress(app)
	p.LoadUI(app)

	messages := make(chan restic.Message)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			var msg restic.Message
			json.Unmarshal(scanner.Bytes(), &msg)
			msg.Body = scanner.Text()
			messages <- msg
		}
	}()

	go p.ListenForMessages(messages)
	app.Run()
}
