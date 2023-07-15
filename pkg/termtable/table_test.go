package termtable

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubeFormat(t *testing.T) {
	data := [][]string{
		{"1/1/2014", "jan_hosting", "2233", "$10.98"},
		{"1/1/2014", "feb_hosting", "2233", "$54.95"},
		{"1/4/2014", "feb_extra_bandwidth", "2233", "$51.00"},
		{"1/4/2014", "mar_hosting", "2233", "$30.00"},
	}

	var buf bytes.Buffer
	OutputTo([]string{"Date", "Description", "CV2", "Amount"}, data, &buf)
	want := `DATE    	DESCRIPTION        	CV2 	AMOUNT 
1/1/2014	jan_hosting        	2233	$10.98	
1/1/2014	feb_hosting        	2233	$54.95	
1/4/2014	feb_extra_bandwidth	2233	$51.00	
1/4/2014	mar_hosting        	2233	$30.00	
`
	assert.Equal(t, want, buf.String(), "kube format rendering failed")
}
