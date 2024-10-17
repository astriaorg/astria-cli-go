module github.com/astriaorg/astria-cli-go/modules/cli

go 1.22.6

require (
	buf.build/gen/go/astria/primitives/protocolbuffers/go v1.35.1-00000000000000-68c4cd4acb9f.1
	buf.build/gen/go/astria/protocol-apis/protocolbuffers/go v1.35.1-00000000000000-42a0beedc983.1
	github.com/99designs/keyring v1.2.2
	github.com/astriaorg/astria-cli-go/modules/bech32m v0.0.0-00010101000000-000000000000
	github.com/astriaorg/astria-cli-go/modules/go-sequencer-client v0.0.0-00010101000000-000000000000
	github.com/gdamore/tcell/v2 v2.7.4
	github.com/pelletier/go-toml/v2 v2.2.2
	github.com/pterm/pterm v0.12.79
	github.com/rivo/tview v0.0.0-20240807205129-e4c497cc59ed
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.1
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	golang.org/x/crypto v0.26.0
)

require (
	github.com/cockroachdb/fifo v0.0.0-20240606204812-0bbfbd93a7ce // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
)

replace github.com/astriaorg/astria-cli-go/modules/go-sequencer-client => ../go-sequencer-client

replace github.com/astriaorg/astria-cli-go/modules/bech32m => ../bech32m

require (
	atomicgo.dev/cursor v0.2.0 // indirect
	atomicgo.dev/keyboard v0.2.9 // indirect
	atomicgo.dev/schedule v0.1.0 // indirect; indirectbuf.build/gen/go/astria/primitives/protocolbuffers/go v1.34.1-20240517150107-e54d335d6779.1 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.4 // indirect
	github.com/btcsuite/btcd/btcutil v1.1.6 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.1 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cometbft/cometbft v0.38.12
	github.com/cometbft/cometbft-db v0.12.0 // indirect
	github.com/containerd/console v1.0.4 // indirect
	github.com/cosmos/gogoproto v1.7.0 // indirect
	github.com/danieljoos/wincred v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dgraph-io/badger/v4 v4.2.0 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.7.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/flatbuffers v24.3.25+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/linxGnu/grocksdb v1.9.1 // indirect
	github.com/lithammer/fuzzysearch v1.1.8 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/petermattis/goid v0.0.0-20240607163614-bb94eb51e7a7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.20.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.etcd.io/bbolt v1.4.0-alpha.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240604190554-fc45aab8b7f8 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/term v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240610135401-a8a62080eff3 // indirect
	google.golang.org/grpc v1.64.1 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
