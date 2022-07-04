
# hacky fix for go customising html coverage output
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
sed -i \
	-e 's/body {/body{ font-family: "Go mono";/g' \
	-e 's/font-weight: bold;//g' \
	-e 's/background: black/background: #DDD/g' \
	-e 's/color: rgb(80, 80, 80);/color: #323232;/g' \
	-e 's/.cov8 { color: rgb(44, 212, 149) }/.cov8 { color: #070; }/g' \
	coverage.html;
sensible-browser coverage.html </dev/null >/dev/null 2>&1 & disown


exit
oriignal="
			body {
				background: black;
				color: rgb(80, 80, 80);
			}
			body, pre, #legend span {
				font-family: Menlo, monospace;
				font-weight: bold;
			}
			#topbar {
				background: black;
				position: fixed;
				top: 0; left: 0; right: 0;
				height: 42px;
				border-bottom: 1px solid rgb(80, 80, 80);
			}
			#content {
				margin-top: 50px;
			}
			#nav, #legend {
				float: left;
				margin-left: 10px;
			}
			#legend {
				margin-top: 12px;
			}
			#nav {
				margin-top: 10px;
			}
			#legend span {
				margin: 0 5px;
			}
			.cov0 { color: rgb(192, 0, 0) }
.cov1 { color: rgb(128, 128, 128) }
.cov2 { color: rgb(116, 140, 131) }
.cov3 { color: rgb(104, 152, 134) }
.cov4 { color: rgb(92, 164, 137) }
.cov5 { color: rgb(80, 176, 140) }
.cov6 { color: rgb(68, 188, 143) }
.cov7 { color: rgb(56, 200, 146) }
.cov8 { color: rgb(44, 212, 149) }
.cov9 { color: rgb(32, 224, 152) }
.cov10 { color: rgb(20, 236, 155) }
"
