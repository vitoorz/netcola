#!/bin/sh

SCRIPT_NAME=$0
TARGET_DIR_1_LEVEL_ABOVE=types

SCRIPT_DIR=${SCRIPT_NAME%/*}
PROTO_FILE=protocol.txt

function doAwkAnalysis() {
bash << EOF
    cat ${PROTO_FILE} | grep -v "^[ \t]*#" | grep ":" | awk -F"[,;]" '${printFunction} {
                for(i=1; i < NF; i++) {
                split(\$i, kv, ":");
                for(j=1; j<=2; j++) {
                    if(kv[j] ~ /[ \t]+"/) {
                        gsub(/[ \t]+"/, "\"", kv[j]);
                    } else {
                        gsub(/[ \t]+/,  "", kv[j]);
                    }
                }

                key = kv[1]; value = kv[2];
                delete kv;
                if(key == "code") {
                    printACode(codeDes);
                    curCode = value;
                }
                codeDes[key] = value;
            }
        }END{
            printACode(codeDes);
        }'
EOF
}
cd ${SCRIPT_DIR}

#=======================================================
# first: generate msg type const
#=======================================================
PACKAGE=types
MSG_TYPE=gnetmsgtype.go

cat << EOF | tee ${MSG_TYPE}
package $PACKAGE

//Auto generated, do not modify unless you know clearly what you are doing.
type MsgType int32

const (
    MT_HeatBeat          =  MsgType(  123456)
    MT_Blank             =  MsgType(  0)

EOF

printFunction='function printACode(codeDesc){
            action = codeDes["action"]
            code   = codeDesc["code"]
            if(action != "" && code != "") {
                type=sprintf("MT_%s", action);
                printf("    %-20s =  MsgType(%3s)\n", type, codeDesc["code"]);
            } else if (action != "") {
                print("    //WARN: INVALID PROTO MAY EXIST HERE")
            }
            delete codeDesc
        }'

doAwkAnalysis | tee -a ${MSG_TYPE}
echo ")" | tee -a ${MSG_TYPE}

#message type in string form
cat << EOF | tee -a ${MSG_TYPE}

func (mt MsgType)ToString() string {
    return netMsgTypeName[mt]
}

var netMsgTypeName = map[MsgType]string {
    MT_HeatBeat         :  "MT_HeatBeat",
    MT_Blank            :  "MT_Blank",
EOF
printFunction='function printACode(codeDesc){
            action = codeDes["action"]
            code   = codeDesc["code"]
            if(action != "" && code != "") {
                type=sprintf("MT_%s", action);
                printf("    %-20s:  \"%s\",\n", type, type);
            } else if (action != "") {
                print("    //WARN: INVALID PROTO MAY EXIST HERE")
            }
            delete codeDesc
        }'

doAwkAnalysis | tee -a ${MSG_TYPE}
echo "}" | tee -a ${MSG_TYPE}


mv ${MSG_TYPE} ../${TARGET_DIR_1_LEVEL_ABOVE}/