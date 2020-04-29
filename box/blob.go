// Code generated by go generate; DO NOT EDIT.
// generated using files from resources directory
package box

import (
	"github.com/filecoin-project/specs-actors/actors/abi"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func init() {
	resources.Add("/TestChainValidationMessageSuiteTestAccountActorCreationfailcreateBLSaccountactorinsufficientbalance", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(180), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacedd7mla4cv6ffsyfn2lwac5mllmrbxflcatd4c6zz4gyyun7ryfl6"}})
	resources.Add("/TestChainValidationMessageSuiteTestAccountActorCreationfailcreateSECP256K1accountactorinsufficientbalance", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(122), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacedd7mla4cv6ffsyfn2lwac5mllmrbxflcatd4c6zz4gyyun7ryfl6"}})
	resources.Add("/TestChainValidationMessageSuiteTestAccountActorCreationsuccesscreateBLSaccountactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 874}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(874), Root: "bafy2bzacebdwqq5jmu5vonozztvy4yxfeipk5biq4e7ewnyhxzbzpp3mkkyti"}})
	resources.Add("/TestChainValidationMessageSuiteTestAccountActorCreationsuccesscreateSECP256K1accountactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 758}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(758), Root: "bafy2bzacebamoezcyi3gna22ku66ot2ym4jbdwyyf4m7bdi67tspt4lk6hah6"}})
	resources.Add("/TestChainValidationMessageSuiteTestInitActorSequentialIDAddressCreate", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacedtds4qjiier2e5mfy7qf6ijekz3u4tlgii27ranp654rn7b6tkwg"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x6a, 0x55, 0x2, 0x24, 0xe3, 0x7, 0xb3, 0xd6, 0x68, 0x4b, 0xbe, 0xc1, 0xc1, 0x1f, 0x1b, 0x80, 0xa8, 0xaf, 0xe0, 0x4, 0x1d, 0xea, 0x94}, GasUsed: 2007}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(2007), Root: "bafy2bzaceaa527u7ehuwuqagycrcbbiqs46ijigm75x2wrpjzjtmgbr43txoe"}})
	resources.Add("/TestChainValidationMessageSuiteTestMessageApplicationEdgecasesabortduringactorexecution", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacedtds4qjiier2e5mfy7qf6ijekz3u4tlgii27ranp654rn7b6tkwg"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 16, ReturnValue: []uint8{}, GasUsed: 349}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(349), Root: "bafy2bzacecvgtc6notbx6yy22jfkj5557zqj6iqqxiahbf33ie4b46wuw7i3y"}})
	resources.Add("/TestChainValidationMessageSuiteTestMessageApplicationEdgecasesfailtocovergascostformessagereceiptonchain", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(112), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaced2nj5y3hh2erhvs2u4ytfuvytanlreapxs52bl64bwvxs3arpb22"}})
	resources.Add("/TestChainValidationMessageSuiteTestMessageApplicationEdgecasesinvalidactorCallSeqNum", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaced2nj5y3hh2erhvs2u4ytfuvytanlreapxs52bl64bwvxs3arpb22"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(88), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaced2nj5y3hh2erhvs2u4ytfuvytanlreapxs52bl64bwvxs3arpb22"}})
	resources.Add("/TestChainValidationMessageSuiteTestMessageApplicationEdgecasesnotenoughgastopaymessageonchainsizecost", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(1120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaced2nj5y3hh2erhvs2u4ytfuvytanlreapxs52bl64bwvxs3arpb22"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(800), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaced2nj5y3hh2erhvs2u4ytfuvytanlreapxs52bl64bwvxs3arpb22"}})
	resources.Add("/TestChainValidationMessageSuiteTestMultiSigActoraddsigner", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1944}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1944), Root: "bafy2bzacebkqojwho7xin7kriszmrs3huymuxgwzyveb4d7tjekejqsmpcx56"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 8, ReturnValue: []uint8{}, GasUsed: 108}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(108), Root: "bafy2bzacec6bbvlo5thdaaennbbd6zg2cp6lhynni4ixcazlhmmvl6xynqo36"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1214}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1214), Root: "bafy2bzacec75urc6j3m65ulibv6yz2yxq4dypvgyg7tuczwi4vggnjzvc4d2w"}})
	resources.Add("/TestChainValidationMessageSuiteTestMultiSigActorconstructortest", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1853}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1853), Root: "bafy2bzacea6ojv5sgnxqcdv62xjxrzyxs3f7jzc3byc22kgbmp5fh4caed67k"}})
	resources.Add("/TestChainValidationMessageSuiteTestMultiSigActorproposeandapprove", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x6a, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 2039}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(2039), Root: "bafy2bzacechfmdh5zwribedzuqrj25vxi7xapb2wi3rc66tyhrowkswovql3o"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 897}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(897), Root: "bafy2bzacea3ewjvploo6jeatxtgoeqqthjxru3o64zhsazfjdbztzf3dn7ngs"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 18, ReturnValue: []uint8{}, GasUsed: 224}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(224), Root: "bafy2bzaceb767hcykmm4j7ftkznsinbgsfdvklqbvbanqwggpgqnzug76af4u"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 18, ReturnValue: []uint8{}, GasUsed: 240}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(240), Root: "bafy2bzacebkerdqlimmpeyfpmoviivw5sj7f2tmwz5bzgoawtb23zqmigevic"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 1166}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1166), Root: "bafy2bzacea6pike2ymttxq7cvrei23ahdf7u5kaag5klhlnijdg43ouzrkmi6"}})
	resources.Add("/TestChainValidationMessageSuiteTestMultiSigActorproposeandcancel", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x6a, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 2039}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(2039), Root: "bafy2bzacedr4l7fp62iyqyrgg6dvpimh5bjp72eamruo2q3acpx6qedsq4du6"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 897}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(897), Root: "bafy2bzacecmoqw56uzjvswkk5itbltyt3h5zh6yj46dpzb4qikvl2doyfltbw"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 18, ReturnValue: []uint8{}, GasUsed: 294}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(294), Root: "bafy2bzacebhwgahdsljw7oar3ywhdltwkwwnf6tml6kz7mbwsfzgazwwdmci6"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 577}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(577), Root: "bafy2bzacecn66aupzbj2dr64v6ajtpkxfsajmrpjp4lfx77o5wejojublyzso"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsfailabortedexec", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 2708}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(2708), Root: "bafy2bzaceahljz4ileev3flhkcz5jp55f7fvxa26jjjqfm5wcduwuxxcvhtqm"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsfailinnerabort", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1017}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1017), Root: "bafy2bzaceaogz4il3ga6a47apto642fq2hsct3uuf2esjvillty4vlny2xwxa"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsfailinvalidmethodnumforactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 954}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(954), Root: "bafy2bzacedokdnqx45k6yqdpeushtubwvgra7q7wjx4oyqsv7r43jf5sbvzje"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsfailinvalidmethodnumnewactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1753}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1753), Root: "bafy2bzaceaabbppf3cf2qdkdkd4ylrzwyuhxbmc5rvi5bmgzlhbwknfapjtym"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsfailmismatchedparams", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1008}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1008), Root: "bafy2bzaced3k7ej4tby63skqxz3cy2764tdjcnj4wweavtly4fe6oym2jjoeg"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsfailmissingparams", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 945}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(945), Root: "bafy2bzaceazrr5iafymea24tvoamluycq3wyeweyvfhhjaoslsxg5ilmz45aq"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsfailnonexistentIDaddress", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 944}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(944), Root: "bafy2bzacebxlxtjwixgwcpxin3jr23q2rfy63eo43nhot5aynkzfuranpo3vm"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsfailnonexistentactoraddress", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1108}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1108), Root: "bafy2bzacebslq3bordmotne6ijgfajky4uujcioqingj7sfz6lwfapfyri53g"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsokbasic", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 935}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(935), Root: "bafy2bzaceb2fzwscd5ldskybvaflhytmjojzbk37zdzhsnn2ipkegjyzeu6ym"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsoknonCBORparamswithtransfer", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1770}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1770), Root: "bafy2bzacedmfrwyrhlypavdnb74p4zsv6psr4mcec5ez4xqxo5fxh4ay6no46"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsokrecursive", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1176}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1176), Root: "bafy2bzaceb7okljy65prqmiiaa7yxxvp6awsilnrom25y6muq75ih7g6invhe"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsoktonewactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1734}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1734), Root: "bafy2bzacedotxkkx4y5mzbjrktque76ce5yxpki27s73lgacylmqec2zmhkse"}})
	resources.Add("/TestChainValidationMessageSuiteTestNestedSendsoktonewactorwithinvoke", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzacebzenpy3tb4bbk574stitg5ag3bjdk2munikuae4i5da2lpujjlmk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1777}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1777), Root: "bafy2bzaceags5ia2dpmsnjv5mdvcjpvpibdfhmz2b7qzk6ukthdddhhazn4ky"}})
	resources.Add("/TestChainValidationMessageSuiteTestPaychhappypathcollect", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacedtds4qjiier2e5mfy7qf6ijekz3u4tlgii27ranp654rn7b6tkwg"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 322}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(322), Root: "bafy2bzacec4rvl2am374murzbj4d2wh4i5dim3g5phyzcghdgatb3zici6pt6"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 191}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(191), Root: "bafy2bzacedhonzbtx3kwhh3mrlgcho52nvd2ugplymx3ww7oql2ggj5vcemvq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 236}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(236), Root: "bafy2bzaceca7w33p7mkwfwk2shfm66y5mwzbtqr36luzkeewfxbi6baskeoru"}})
	resources.Add("/TestChainValidationMessageSuiteTestPaychhappypathconstructor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacedtds4qjiier2e5mfy7qf6ijekz3u4tlgii27ranp654rn7b6tkwg"}})
	resources.Add("/TestChainValidationMessageSuiteTestPaychhappypathupdate", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacedtds4qjiier2e5mfy7qf6ijekz3u4tlgii27ranp654rn7b6tkwg"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 320}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(320), Root: "bafy2bzaceamqf3xbvjcbqf5yxjicrwyssababrrdjvswnwf4zpnwrrtqx4g3q"}})
	resources.Add("/TestChainValidationMessageSuiteTestValueTransferAdvancefailtotransferfromunknownaccounttoknownaddress", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaced2nj5y3hh2erhvs2u4ytfuvytanlreapxs52bl64bwvxs3arpb22"}})
	resources.Add("/TestChainValidationMessageSuiteTestValueTransferAdvancefailtotransferfromunknownaddresstounknownaddress", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacec4kgq2yfzhe2lyinmkkhqtm3weyzt3fvsubn7ttbqhb4p5r5v6es"}})
	resources.Add("/TestChainValidationMessageSuiteTestValueTransferAdvanceoktransferfromknownaddresstonewaccount", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 756}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(756), Root: "bafy2bzacedc77zu3n62rtp5aqqa6fiwo5t22otbadpwflhx5aycmq62i7wvz6"}})
	resources.Add("/TestChainValidationMessageSuiteTestValueTransferAdvanceselftransfer", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(130), Root: "bafy2bzaceaei54fghnrkhctifwlafddfrxttrbeytbns6yo5kcvxgya375jzm"}})
	resources.Add("/TestChainValidationMessageSuiteTestValueTransferSimplefailtotransfermorefundsthansenderbalance0", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(124), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaceb5h7lia3dfihdwazfoapzcdl73usdyryj6emyxnn7vfn5ichgsiq"}})
	resources.Add("/TestChainValidationMessageSuiteTestValueTransferSimplefailtotransfermorefundsthansenderhaswhensenderbalancezero", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacea72opjrl4s2i6wqhqbpbjliwh32rpdliinvwnqs2leggweecy7to"}})
	resources.Add("/TestChainValidationMessageSuiteTestValueTransferSimplesuccessfullytransferfundsfromsendertoreceiver", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(130), Root: "bafy2bzaceagklcxqn3sgf2kceg5rp7v23qytut7pxlo5dqvom3ghuib537evi"}})
	resources.Add("/TestChainValidationMessageSuiteTestValueTransferSimplesuccessfullytransferzerofundsfromsendertoreceiver", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 126}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(126), Root: "bafy2bzaceaasln7veg7xgfxnjvphiznmp5effnokfu7otd4qja3k7nxyzheya"}})
	resources.Add("/TestChainValidationTipSetSuiteTestBlockMessageDeduplicationapplyaduplicatedBLSmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}}, Root: "bafy2bzacea67ulf5c3pflo7mxrrb4rmtxd76vwvlysiq6tv527w3d3p2zr35u"}})
	resources.Add("/TestChainValidationTipSetSuiteTestBlockMessageDeduplicationapplyasingleBLSmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}}, Root: "bafy2bzacea67ulf5c3pflo7mxrrb4rmtxd76vwvlysiq6tv527w3d3p2zr35u"}})
	resources.Add("/TestChainValidationTipSetSuiteTestBlockMessageDeduplicationapplyasingleSECPmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 268}}, Root: "bafy2bzacecgilzzfh4b3bqlzkx2jft4pd77qleljjk2v4nr6egiq7pvrx6cv2"}})
	resources.Add("/TestChainValidationTipSetSuiteTestBlockMessageDeduplicationapplyduplicateSECPmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 268}}, Root: "bafy2bzacecgilzzfh4b3bqlzkx2jft4pd77qleljjk2v4nr6egiq7pvrx6cv2"}})
	resources.Add("/TestChainValidationTipSetSuiteTestMinerRewardsAndPenaltiesminerpenaltyexceedsdeclaredgaslimitforBLSmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzaceb2ahxu7fzhd7zdmmnrmx2vt2lkgqsa3fomt6mgh2rrg4pijxudkg"}})
	resources.Add("/TestChainValidationTipSetSuiteTestMinerRewardsAndPenaltiesminerpenaltyexceedsdeclaredgaslimitforSECPmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzacebhqzgrl2dgloxh7ixss47zkwrvrw363cctzh3ijt7nbevet7yres"}})
	resources.Add("/TestChainValidationTipSetSuiteTestMinerRewardsAndPenaltiesoksimplesend", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}, types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}}, Root: "bafy2bzaceaynvfji76eob2s22im75yfdaeji7t6lpkslyqord2s3q4nixmyyw"}, types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 92}, types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 92}}, Root: "bafy2bzaceck62r3iuo47qludxgw5bfxjeehdtgepgqzqsscmn2y55qfh5oiki"}, types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 92}, types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 92}}, Root: "bafy2bzaceb5yn7k6fo7wjarjy7n3pjkkirypsmjpww2aif4spifpdsozvk6es"}, types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 54}, types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 54}}, Root: "bafy2bzacebqcsh7ki3mbo5vxfgndpyedtd4atny6wajwlyafhkpx3h4f3ryck"}})
	resources.Add("/TestChainValidationTipSetSuiteTestMinerRewardsAndPenaltiespenalizesenderdoestexist", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzacedrrdlticotdl6qnww2fox4xqrvl44haaksn7z7dakpxjgqgqutwe"}})
	resources.Add("/TestChainValidationTipSetSuiteTestMinerRewardsAndPenaltiespenalizesenderinsufficientbalance", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 58}, types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzaceduzqqrtpsyagfqrv3xuynz4z6rt7supc5cmtyribv4lboma5glcq"}})
	resources.Add("/TestChainValidationTipSetSuiteTestMinerRewardsAndPenaltiespenalizesendernonaccount", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzacedernbi3nhvmdxnlsekcpzmuegom33hl74o77kc4sfnotyxhz4k7s"}})
	resources.Add("/TestChainValidationTipSetSuiteTestMinerRewardsAndPenaltiespenalizewrongcallseqnum", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzaceaomweyl3dl6gzkqxmoqebwtjqacamwterogcdovd4rxj4pqdbkk4"}})
}
