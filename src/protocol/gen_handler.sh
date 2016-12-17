#!/bin/sh

SCRIPT_NAME=$0
SCRIPT_DIR=${SCRIPT_NAME%/*}
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

printFunctionMap='function printACode(codeDesc){
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

printFunctionOn='function printACode(codeDesc){
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

            printf("\nfunc On_%s(objectId IdString, opCode MsgType, payLoad []byte) interface{} {\n",action);
            printf("\treq := &%s{}\n\tpb.Unmarshal(payLoad, req)\n", payload);
            printf("\treturn Handle_%s(objectId, opCode, req)\n}\n", action);
            delete codeDesc
        }'

 printFunctionHandleReq='function printACode(codeDesc){
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
                printf("\nfunc Handle_%s(objectId IdString, opCode MsgType, req *%s) interface{} {\n",action, payload);
                printf("\treturn nil\n}\n");
            }

            delete codeDesc
        }'

 printFunctionHandleAck='function printACode(codeDesc){
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

            if(codeDesc["reaction"] == "") {
                printf("\nfunc Handle_%s(objectId IdString, opCode MsgType, ack *%s) interface{} {\n",action, payload);
                printf("\treturn nil\n}\n");
            }

            delete codeDesc
        }'
#=======================================================
#BEGIN TO PROCESS FILE
#=======================================================
todo_list=`ls *protocol.txt`
for PROTO_FILE in ${todo_list}; do
    #=======================================================
    # generate msg handler map for net
    #=======================================================
    TARGET_DIR_1_LEVEL_ABOVE=`cat ${PROTO_FILE} | grep "package=" | awk -F= '{print $2}'`
    OUTPUT_FILE=handlermap.go
    PACKAGE=${TARGET_DIR_1_LEVEL_ABOVE##*/}

cat << EOF | tee ${OUTPUT_FILE}
package ${PACKAGE}
//Auto generated, do not modify unless you know clearly what you are doing.
import . "types"

type NetMsgHandler func(objectId IdString, opcode MsgType, content []byte ) interface{};

type NetMsgCb struct{
	OpCode  MsgType
	RetCode MsgType
	Handler NetMsgHandler
	Desc    string
}

var NetMsgTypeHandler = map[MsgType]*NetMsgCb {
EOF

    printFunction=${printFunctionMap}
    doAwkAnalysis | tee -a ${OUTPUT_FILE}
    echo "}" | tee -a ${OUTPUT_FILE}
    mv ${OUTPUT_FILE} ../${TARGET_DIR_1_LEVEL_ABOVE}/

    #=======================================================
    # second: generate msg handler for each net payload
    #=======================================================
    OUTPUT_FILE=handle.go

cat << EOF | tee ${OUTPUT_FILE}
package ${PACKAGE}
//Auto generated, do not modify unless you know clearly what you are doing.

import (
	pb "github.com/golang/protobuf/proto"
	. "types"
EOF
    cat ${PROTO_FILE} | grep "gamedealer" | awk -F= '{split($2, ps, " "); for(idx in ps) printf("    . \"%s\"\n"), ps[idx]}'| tee -a ${OUTPUT_FILE}
    echo ")" | tee -a ${OUTPUT_FILE}

    printFunction=${printFunctionOn}
    doAwkAnalysis | tee -a ${OUTPUT_FILE}
    mv ${OUTPUT_FILE} ../${TARGET_DIR_1_LEVEL_ABOVE}/

    #=======================================================
    # third: generate msg handler for each req payload
    #=======================================================
    OUTPUT_FILE=handlereq.go

cat << EOF | tee ${OUTPUT_FILE}
package ${PACKAGE}
//Auto generated, do not modify unless you know clearly what you are doing.
import . "types"
EOF
    printFunction=${printFunctionHandleReq}
    doAwkAnalysis | tee -a ${OUTPUT_FILE}
    echo "mv ${OUTPUT_FILE} ../${TARGET_DIR_1_LEVEL_ABOVE}/${OUTPUT_FILE}.auto"
    mv ${OUTPUT_FILE} ../${TARGET_DIR_1_LEVEL_ABOVE}/${OUTPUT_FILE}.auto

    #=======================================================
    # fourth: generate msg handler for each ack payload
    #=======================================================
    OUTPUT_FILE=handleack.go

cat << EOF | tee ${OUTPUT_FILE}
package ${PACKAGE}
//Auto generated, do not modify unless you know clearly what you are doing.
import . "types"
EOF
    printFunction=${printFunctionHandleAck}
    doAwkAnalysis | tee -a ${OUTPUT_FILE}
    echo "mv ${OUTPUT_FILE} ../${TARGET_DIR_1_LEVEL_ABOVE}/${OUTPUT_FILE}.auto"
    mv ${OUTPUT_FILE} ../${TARGET_DIR_1_LEVEL_ABOVE}/${OUTPUT_FILE}.auto
done