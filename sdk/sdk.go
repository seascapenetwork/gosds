/*The gosds/sdk package is the client package to interact with SDS.
The following commands are available in this SDK:

1. Subscribe - subscribe for events
2. Sign - send a transaction to the blockchain
3. AddToPool - send a transaction to the pool that will be broadcasted to the blockchain bundled.

*/
package sdk

import (
	"fmt"
)

var Version string = "1.0.0"

func PrintVersion() {
	fmt.Println("[gosds] " + Version)
}
