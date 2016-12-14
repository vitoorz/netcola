#!/bin/sh

SCRIPT_NAME=$0
TARGET_DIR_1_LEVEL_ABOVE=types
SCRIPT_DIR=${SCRIPT_NAME%/*}
PACKAGE=types
cd ${SCRIPT_DIR}

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

#=======================begin===========================
todo_list=`ls *protocol.txt`
for PROTO_FILE in ${todo_list}; do
    #=======================================================
    # first: generate msg type const
    #=======================================================
    MSG_TYPE=`cat ${PROTO_FILE} | grep "go="| awk -F= '{print $2}'`

    cat << EOF | tee ${MSG_TYPE}
package ${PACKAGE}

//Auto generated, do not modify unless you know clearly what you are doing.

const (
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

var ${MSG_TYPE%.go}Names = map[MsgType]string {
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

    echo "mv ${MSG_TYPE} ../${TARGET_DIR_1_LEVEL_ABOVE}/"
    mv ${MSG_TYPE} ../${TARGET_DIR_1_LEVEL_ABOVE}/

done