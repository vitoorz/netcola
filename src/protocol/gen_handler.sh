#!/bin/sh

SCRIPT_NAME=$0
TARGET_DIR_1_LEVEL_ABOVE=play

SCRIPT_DIR=${SCRIPT_NAME%/*}
PROTO_FILE=protocol.txt

OUTPUT_FILE=

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

#=======================================================
#BEGIN TO PROCESS FILE
#=======================================================

#=======================================================
# second: generate msg handler map for net
#=======================================================
cd ${SCRIPT_DIR}
OUTPUT_FILE=handlermap.go
PACKAGE=${TARGET_DIR_1_LEVEL_ABOVE}

cat << EOF | tee ${OUTPUT_FILE}
package ${PACKAGE}
//Auto generated, do not modify unless you know clearly what you are doing.
import . "types"

type NetMsgHandler func(playerId IdString, opcode MsgType, content []byte ) interface{};

type NetMsgCb struct{
	OpCode  MsgType
	RetCode MsgType
	Handler NetMsgHandler
	Desc    string
}

var NetMsgTypeHandler = map[MsgType]*NetMsgCb {
EOF

printFunction='function printACode(codeDesc){
            action = codeDes["action"]
            if(action == "") {
                if (codeDes["code"] != "") print("    //WARN: INVALID PROTO MAY EXIST HERE")
                delete codeDesc
                return
            }

            desc = codeDesc["desc"]
            if( desc == "") {
                desc = "\"\""
            }
            reaction = codeDesc["reaction"]
            if( reaction == "") {
                reaction = "Blank"
            }

            opcode=sprintf("MT_%s", action);
            retcode=sprintf("MT_%s", reaction);
            handler=sprintf("On_%s", action)
            printf("\t%-16s :&NetMsgCb{%-16s, %-16s, %-16s, %s},\n",
                  opcode, opcode, retcode, handler, desc);

            delete codeDesc
        }'

doAwkAnalysis | tee -a ${OUTPUT_FILE}
echo "}" | tee -a ${OUTPUT_FILE}
mv ${OUTPUT_FILE} ../${TARGET_DIR_1_LEVEL_ABOVE}/

#=======================================================
# third: generate msg handler for each net payload
#=======================================================
OUTPUT_FILE=handle.go

cat << EOF | tee ${OUTPUT_FILE}
package ${PACKAGE}
//Auto generated, do not modify unless you know clearly what you are doing.

import (
	pb "github.com/golang/protobuf/proto"
	. "types"
)

EOF

printFunction='function printACode(codeDesc){
            action = codeDes["action"]
            if(action == "") {
                if (codeDes["code"] != "") print("    //WARN: INVALID PROTO MAY EXIST HERE")
                delete codeDesc
                return
            }

            reaction = codeDesc["reaction"]
            if( reaction == "") {
                reaction = "Blank"
            }

            payload = codeDesc["payload"]
            if(payload  == "") {
                payload = action
            }

            if(reaction == "Blank") {
                printf("\nfunc On_%s(playerId IdString, opCode MsgType, payLoad []byte) interface{} {\n",action);
                printf("\treturn Play_InvalidReq(playerId, opCode, payLoad)\n}\n");
            } else {
                printf("\nfunc On_%s(playerId IdString, opCode MsgType, payLoad []byte) interface{} {\n",action);
                printf("\treq :=&%s{}\n\tpb.Unmarshal(payLoad, req)\n", payload);
                printf("\treturn Play_%s(playerId, opCode, req)\n}\n", action);
            }
            delete codeDesc
        }'

doAwkAnalysis | tee -a ${OUTPUT_FILE}
mv ${OUTPUT_FILE} ../${TARGET_DIR_1_LEVEL_ABOVE}/

#=======================================================
# fourth: generate msg handler for each net payload
#=======================================================
OUTPUT_FILE=handleplay.go

cat << EOF | tee ${OUTPUT_FILE}
package ${PACKAGE}
//Auto generated, do not modify unless you know clearly what you are doing.

import . "types"

func Play_InvalidReq(playerId IdString, opCode MsgType, req interface{}) interface{} {
	return getCommonAck(ERR_INVALID_REQ)
}
EOF

printFunction='function printACode(codeDesc){
           action = codeDes["action"]
            if(action  == "") {
                if (codeDes["action"] != "") print("    //WARN: INVALID PROTO MAY EXIST HERE")
                delete codeDesc
                return
            }

            payload = codeDesc["payload"]
            if(payload  == "") {
                payload = action
            }

            if(codeDesc["reaction"] != "") {
                printf("\nfunc Play_%s(playerId IdString, opCode MsgType, req *%s) interface{} {\n",action, payload);
                printf("\treturn nil\n}\n");
            }

            delete codeDesc
        }'

doAwkAnalysis | tee -a ${OUTPUT_FILE}
mv ${OUTPUT_FILE} ../${TARGET_DIR_1_LEVEL_ABOVE}/${OUTPUT_FILE}.auto