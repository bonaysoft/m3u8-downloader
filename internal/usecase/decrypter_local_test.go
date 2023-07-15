package usecase

import (
	"testing"

	"github.com/bonaysoft/m3u8-downloader/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	d, err := NewDecrypterLocal(entity.Properties{Cmd: []string{"../../.local/xet/main.js"}})
	assert.NoError(t, err)

	mu := entity.NewM3u8URL("W$siZGVmaW5pdGlvbl9uYW@lIjoiXHU5YWQ%XHU#ZTA@IiwiZGVmaW5pdGlvbl9wIjoiNzIwUCIsInVybCI6Imh0dHBzOlwvXC9wcmktY#RuLXR%LnhpYW9la#5vdy5jb#@cL#FwcGVpZnB@aml@NzI#OVwvcHJpdmF0ZV9pbmRleFwvMTY#MzM%OTAzNmtHNUhzOC5tM$U%P$NpZ#%9MjVkY#Q%Y#E#YmFlMDg#ODlkZTBmMjM%YTBjZTcwNDUmdD0#NGFhYjFiYSIsImlzX$N@cHBvcnQiOnRydWUsImV%dCI6eyJob$N0IjoiaHR0cHM6XC9cL#J0dC@#b#QueGlhb#Vrbm9$LmNvbSIsInBhdGgiOiIyOTE5ZGY%OHZvZHRyYW5zY$ExMjUyNTI0MTI#XC8wMGI5ZGQ5YjM%NzcwMjMwNTE0MDAwMTc@MVwvZHJtIiwicGFyYW0iOiJzaWduPWRkMWM#MDhiNGI5NTIzN#Q0NzgyOTdhN#QxODdlZTU@JnQ9NjRhYjVhN#EmdXM9S@pUbW@$ZWR$TCJ9fV0=__ba")
	assert.NoError(t, d.M3u8URLDecrypt(mu))
	assert.NotEmpty(t, mu.PlainURL)
	assert.NotEmpty(t, mu.TsURLPart.Host)

	raw := []byte{6, 28, 178, 34, 211, 151, 172, 45, 236, 198, 131, 230, 5, 85, 186, 247}
	decrypted, err := d.KeyDecrypt(raw)
	assert.NoError(t, err)
	assert.Equal(t, []byte{115, 67, 132, 22, 228, 243, 158, 78, 219, 160, 176, 214, 103, 96, 143, 168}, decrypted)
}
