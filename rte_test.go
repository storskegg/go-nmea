package nmea_test

// import (
// 	"testing"

// 	. "github.com/munnik/go-nmea"
// 	"github.com/stretchr/testify/assert"
// )

// var rtetests = []struct {
// 	name string
// 	raw  string
// 	err  string
// 	msg  RTE
// }{
// 	{
// 		name: "good sentence",
// 		raw:  "$IIRTE,4,1,c,Rte 1,411,412,413,414,415*6F",
// 		msg: RTE{
// 			NumberOfSentences:         NewInt64(4),
// 			SentenceNumber:            NewInt64(1),
// 			ActiveRouteOrWaypointList: NewString(ActiveRoute),
// 			Name:                      NewString("Rte 1"),
// 			Idents:                    []String{NewString("411"), NewString("412"), NewString("413"), NewString("414"), NewString("415")},
// 		},
// 	},
// 	{
// 		name: "index out if range",
// 		raw:  "$IIRTE,4,1,c,Rte 1*77",
// 		err:  "nmea: IIRTE invalid ident of waypoints: index out of range",
// 	},
// 	{
// 		name: "invalid number of sentences",
// 		raw:  "$IIRTE,X,1,c,Rte 1,411,412,413,414,415*03",
// 		err:  "nmea: IIRTE invalid number of sentences: X",
// 	},
// 	{
// 		name: "invalid sentence number",
// 		raw:  "$IIRTE,4,X,c,Rte 1,411,412,413,414,415*06",
// 		err:  "nmea: IIRTE invalid sentence number: X",
// 	},
// 	{
// 		name: "invalid active route or waypoint list",
// 		raw:  "$IIRTE,4,1,X,Rte 1,411,412,413,414,415*54",
// 		err:  "nmea: IIRTE invalid active route or waypoint list: X",
// 	},
// }

// func TestRTE(t *testing.T) {
// 	for _, tt := range rtetests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			m, err := Parse(tt.raw)
// 			if tt.err != "" {
// 				assert.Error(t, err)
// 				assert.EqualError(t, err, tt.err)
// 			} else {
// 				assert.NoError(t, err)
// 				rte := m.(RTE)
// 				rte.BaseSentence = BaseSentence{}
// 				assert.Equal(t, tt.msg, rte)
// 			}
// 		})
// 	}
// }
