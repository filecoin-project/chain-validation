// Code generated by go generate; DO NOT EDIT.
// generated using files from resources directory
package box

import (
	"github.com/filecoin-project/specs-actors/actors/abi"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func init() {
	resources.Add("/MessageTestAccountActorCreationfailcreateBLSaccountactorinsufficientbalance", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(180), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacebbwghkqcy3o6yzs75cqegsh6ymwsiczl6sesxicjtdbxgufmqq6q"}})
	resources.Add("/MessageTestAccountActorCreationfailcreateSECP256K1accountactorinsufficientbalance", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(122), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacebbwghkqcy3o6yzs75cqegsh6ymwsiczl6sesxicjtdbxgufmqq6q"}})
	resources.Add("/MessageTestAccountActorCreationsuccesscreateBLSaccountactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 874}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(874), Root: "bafy2bzaceagmq4q2zn3lwqg2ncy2zmldq434cl6x7bdpfuy7egsbj3xhxa57i"}})
	resources.Add("/MessageTestAccountActorCreationsuccesscreateSECP256K1accountactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 758}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(758), Root: "bafy2bzacebf2pgd2bqli25ivxqqoa3rpytlx2ybdfxwtahb5ci27qux6b4bfo"}})
	resources.Add("/MessageTestInitActorSequentialIDAddressCreate", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacecnonppezshaioezjllz5rnxfvxosskbzo4rdw54ckefyxfxhuewc"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x6a, 0x55, 0x2, 0x24, 0xe3, 0x7, 0xb3, 0xd6, 0x68, 0x4b, 0xbe, 0xc1, 0xc1, 0x1f, 0x1b, 0x80, 0xa8, 0xaf, 0xe0, 0x4, 0x1d, 0xea, 0x94}, GasUsed: 2007}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(2007), Root: "bafy2bzaceavzz7yrgjzj664oxgkz6psyrbwtd6lsli2ywt2wuvybarbna4qzi"}})
	resources.Add("/MessageTestMessageApplicationEdgecasesabortduringactorexecution", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacecnonppezshaioezjllz5rnxfvxosskbzo4rdw54ckefyxfxhuewc"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 16, ReturnValue: []uint8{}, GasUsed: 349}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(349), Root: "bafy2bzacebi45dh37ao2wcdgfmvkdcrnzgi6zwfw3owt4hqjr2fg6ul6foqli"}})
	resources.Add("/MessageTestMessageApplicationEdgecasesfailnotenoughgastocoveraccountactorcreation", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 756}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(756), Root: "bafy2bzaceb3lbjrk2uo462tmajh3kvwzli3t5v3fiyzyvci3dhf4bxrnf4bek"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 706}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(706), Root: "bafy2bzacec2i65peh4swtuvj35rz4erecra3lemxvyqtm5lxgj3mre2rcbifg"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 656}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(656), Root: "bafy2bzacebbolo2nd4r6ssz5c374lpmvmu5kma45vellkm5p2sctv5rgja7x2"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 606}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(606), Root: "bafy2bzacearlinj5d2axu2dnhgn23l5hfuhhzioyd2joyhpag4tmh34ubskim"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 556}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(556), Root: "bafy2bzaceank5n3q2quxwp3cdrfmy5qoblabs6gokxltozcfvdzb6bdfcocoa"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 506}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(506), Root: "bafy2bzacecqbkfdbkftl7cd3kp7xpwehvervkpqnzknmylyfkjd55fkdk5ctq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 456}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(456), Root: "bafy2bzacecxidvf7f2odahnp3ql3xgzfnbu47sogveauvxeiet6jsmy4kn4ig"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 406}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(406), Root: "bafy2bzacea4nqyqg7kl46qgzrlzuyqoyg4khzsu5otah7zafui6kxys5rax42"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 356}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(356), Root: "bafy2bzacebdaqzy5d637htwdvfl2banraumd3sesmjj73ptco4manlg4cnass"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 306}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(306), Root: "bafy2bzacebkwrf74ofygsn54d2pbhs6krargu4llffw6kwry26g5jimkvp5xw"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 256}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(256), Root: "bafy2bzaceaqtzixcvkm4fvmgalgbw24or44uewxyjhr2orare32p7yliimilk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 206}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(206), Root: "bafy2bzacedodbzlfi7cedojsmoomwmsyg6d6xfej3kl3k2vrpsw4atd4fx5fo"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 156}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(156), Root: "bafy2bzaceclc67qzxqybgw74zlqqvje4fu24so4crtoj266dqpg6x2fu4arb6"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(114), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaceclc67qzxqybgw74zlqqvje4fu24so4crtoj266dqpg6x2fu4arb6"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(114), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaceclc67qzxqybgw74zlqqvje4fu24so4crtoj266dqpg6x2fu4arb6"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(112), Reward: abi.NewTokenAmount(0), Root: "bafy2bzaceclc67qzxqybgw74zlqqvje4fu24so4crtoj266dqpg6x2fu4arb6"}})
	resources.Add("/MessageTestMessageApplicationEdgecasesfailtocovergascostformessagereceiptonchain", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(112), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacecxzkn6usolqt6tf5pg5lrhh6jqetxzqv5hfsn7aqe2xesyilh4pg"}})
	resources.Add("/MessageTestMessageApplicationEdgecasesinvalidactorCallSeqNum", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacecxzkn6usolqt6tf5pg5lrhh6jqetxzqv5hfsn7aqe2xesyilh4pg"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(88), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacecxzkn6usolqt6tf5pg5lrhh6jqetxzqv5hfsn7aqe2xesyilh4pg"}})
	resources.Add("/MessageTestMessageApplicationEdgecasesinvalidmethodforreceiver", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 3, ReturnValue: []uint8{}, GasUsed: 138}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(138), Root: "bafy2bzaceaab4jglcqdoicp3oiuoc3v4zs2emw56zvbzier7ehgvke3vcybty"}})
	resources.Add("/MessageTestMessageApplicationEdgecasesnotenoughgastopaymessageonchainsizecost", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(1120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacecxzkn6usolqt6tf5pg5lrhh6jqetxzqv5hfsn7aqe2xesyilh4pg"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(800), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacecxzkn6usolqt6tf5pg5lrhh6jqetxzqv5hfsn7aqe2xesyilh4pg"}})
	resources.Add("/MessageTestMessageApplicationEdgecasesreceiverIDActoraddressdoesnotexist", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 5, ReturnValue: []uint8{}, GasUsed: 98}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(98), Root: "bafy2bzaceavsu6wxqkwlhoj2cabktgdteaqrq6zcry3zqg37g6r347s45z75g"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 5, ReturnValue: []uint8{}, GasUsed: 130}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(130), Root: "bafy2bzaceb6vktnzgy4ujbszwiqu4qg36jrun4fyqz4trbjh3ntej5k6q3wbe"}})
	resources.Add("/MessageTestMultiSigActoraddsigner", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1944}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1944), Root: "bafy2bzaceb3htqvsrgbdmkxr6raixwpuiq5kqkwwyecrszhp4jlft2kycuy2a"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 8, ReturnValue: []uint8{}, GasUsed: 108}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(108), Root: "bafy2bzaceddhpwpmha2dwr5joqoja73t3v66xqnggfz77suwz2g7smllxpbi4"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1214}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1214), Root: "bafy2bzaceapg2yijdy7qvkp3cowa2egnffr2kz2gtaeq65cnzj67kzwybs3nm"}})
	resources.Add("/MessageTestMultiSigActorconstructortest", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1853}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1853), Root: "bafy2bzaceal3q2uen4wio2lxtpunwt237autbzgxhhvgrhso5zlwycqkw7yn6"}})
	resources.Add("/MessageTestMultiSigActorproposeandapprove", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x6a, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 2039}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(2039), Root: "bafy2bzacedpa5r33sclgdao3ipo552gulfhmk7owomerhj6e7mtmrodtbzquk"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 897}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(897), Root: "bafy2bzaceatexpnyib6sc7idi7caoslwkjui5wgaoy2b3yzq2wvfqiekj4l3y"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 18, ReturnValue: []uint8{}, GasUsed: 224}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(224), Root: "bafy2bzaceakj6anbq3kaump5apvmgpfvgpywko3samnaw4uzqibhccly3ctfq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 18, ReturnValue: []uint8{}, GasUsed: 240}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(240), Root: "bafy2bzacebpvcc5etr2clcfmeaz5cbj2vf2v2jp4v77ysqbxv6wfctvzdvqhg"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 1166}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1166), Root: "bafy2bzacedwaol7oef2l22qrbqwdyhfhvi4dmnbgncv2dk5vj7nfxtlmgmhwe"}})
	resources.Add("/MessageTestMultiSigActorproposeandcancel", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x6a, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 2039}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(2039), Root: "bafy2bzacedadktukz7infefopjfe5v5kn7qy2zantfrjzwodjd5exjirjahte"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 897}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(897), Root: "bafy2bzacedayxurchkwdkfbiytci672o533lhoo2bvs4auotkff4krmt3rmcq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 18, ReturnValue: []uint8{}, GasUsed: 294}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(294), Root: "bafy2bzacebpr5tbe7rjjcmuvsu6xv7yrsjfuqeo5syznq7qnro7oo3myxlv3q"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 577}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(577), Root: "bafy2bzacedq3mrkokik5k3oaysql3cwwys3dkqszptyo3wleanfnffyru5nec"}})
	resources.Add("/MessageTestNestedSendsfailabortedexec", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 2708}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(2708), Root: "bafy2bzacedm7cjujcmgsmcrqfq5tdirjkq4ddqwmgzjd3tgydn3is776wgxjg"}})
	resources.Add("/MessageTestNestedSendsfailinnerabort", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1017}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1017), Root: "bafy2bzacealx4rzm53kd7m4cbmj44oxk7pueszs537wdp66ubgxegavs4rulo"}})
	resources.Add("/MessageTestNestedSendsfailinsufficientfundsfortransferininnersend", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 166}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(166), Root: "bafy2bzaceaa37adqnu3jdb5bpm6e4lb3ug7mahn7rneqngqj7nf65dixqu5s4"}})
	resources.Add("/MessageTestNestedSendsfailinvalidmethodnumforactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 954}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(954), Root: "bafy2bzaceayqu4yoldmsdynkyg7novj7iwieo3chorsd42lkrcokfovg4hoao"}})
	resources.Add("/MessageTestNestedSendsfailinvalidmethodnumnewactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1753}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1753), Root: "bafy2bzacea7f6t6ilbbyvpc5csefwhsn73e6lp5m7e6rmzsql3di5leeqigla"}})
	resources.Add("/MessageTestNestedSendsfailmismatchedparams", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1008}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1008), Root: "bafy2bzacedcgkpn4p2r234kr4iyac6luhbxfi7urrrbez5j3f7h7yd3kl7s2i"}})
	resources.Add("/MessageTestNestedSendsfailmissingparams", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 945}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(945), Root: "bafy2bzaceacbifqmkq72n35if655jm7xagbik763sdjha3cs7k5icp2yen4xi"}})
	resources.Add("/MessageTestNestedSendsfailnonexistentIDaddress", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 944}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(944), Root: "bafy2bzaceaffehpj5gp6bdcud4yvdrhdqathe3rhq4r3xb23vp3ehuujtidms"}})
	resources.Add("/MessageTestNestedSendsfailnonexistentactoraddress", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1108}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1108), Root: "bafy2bzacecluvzfn3zql5gzh6qtfp7hwq7zqtzryz43v63mnyqmjo4m5rt2sk"}})
	resources.Add("/MessageTestNestedSendsokbasic", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 935}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(935), Root: "bafy2bzacedxgrickkargisewyost6pdhzvtfgbayk6gxvo5cf5zvdkwkhgdzy"}})
	resources.Add("/MessageTestNestedSendsoknonCBORparamswithtransfer", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1770}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1770), Root: "bafy2bzacebavjq7blw4k7gtsrr7ol6ax3zu2nmav3wip2bvcxu57dgaocqykk"}})
	resources.Add("/MessageTestNestedSendsokrecursive", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1176}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1176), Root: "bafy2bzacedam5m4nbyspuvpurz2qlcby64z4tum23tgeij4myjsjzojve2ili"}})
	resources.Add("/MessageTestNestedSendsoktonewactor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1734}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1734), Root: "bafy2bzaceb5flj6oivekmdaj2wammnl63z3v4rlsqwedeboeveopx7dt2cl4y"}})
	resources.Add("/MessageTestNestedSendsoktonewactorwithinvoke", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x68, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1815}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1815), Root: "bafy2bzaceax4y3fzk3eovi32c5yqya34tgmcwqv3yrlvonn2nzf66i3edwhlq"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x0}, GasUsed: 1777}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1777), Root: "bafy2bzacebhwjisx3e2z4ka4efeek7fbp26bgf7v4ne6b34nr2v7pr3edn75w"}})
	resources.Add("/MessageTestPaychhappypathcollect", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacecnonppezshaioezjllz5rnxfvxosskbzo4rdw54ckefyxfxhuewc"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 322}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(322), Root: "bafy2bzacedqocp37xa2gpwxa7b24n4jhecdrfd4y6xh6oqz5oqmngir2olnzm"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 191}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(191), Root: "bafy2bzacecvdjaffyidca25f4rgl5ixl7vfddqqh7u5vfsy7n6ns22l5drxe4"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 236}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(236), Root: "bafy2bzacecxsv7hsx67uww4peekyr23hhbhepjm6sh5qlvzgvhcfyyy2sjxcc"}})
	resources.Add("/MessageTestPaychhappypathconstructor", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacecnonppezshaioezjllz5rnxfvxosskbzo4rdw54ckefyxfxhuewc"}})
	resources.Add("/MessageTestPaychhappypathupdate", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x69, 0x55, 0x2, 0x6b, 0x15, 0x6c, 0xef, 0x87, 0xc2, 0x1, 0x9d, 0x65, 0x2f, 0x92, 0x67, 0x83, 0x6, 0x36, 0xb4, 0xd2, 0x4f, 0xac, 0xe1}, GasUsed: 1920}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(1920), Root: "bafy2bzacecnonppezshaioezjllz5rnxfvxosskbzo4rdw54ckefyxfxhuewc"}, types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 320}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(320), Root: "bafy2bzacedanvaelbfxdpgczqq543j2evpkbnkpvxdsu3nvt3kf5h3cke2faq"}})
	resources.Add("/MessageTestValueTransferAdvancefailtotransferfromunknownaccounttoknownaddress", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacecxzkn6usolqt6tf5pg5lrhh6jqetxzqv5hfsn7aqe2xesyilh4pg"}})
	resources.Add("/MessageTestValueTransferAdvancefailtotransferfromunknownaddresstounknownaddress", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacear3f5uexshspifu3bsbzs737uudh4xh4ruq5i6ox2ctc6iuw64ty"}})
	resources.Add("/MessageTestValueTransferAdvanceoktransferfromknownaddresstonewaccount", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 756}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(756), Root: "bafy2bzacec7wn6mvsyybvt3prkhvn6hemqmcdutl2u2bj2nnjvmryispa4agy"}})
	resources.Add("/MessageTestValueTransferAdvanceselftransfer", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(130), Root: "bafy2bzacedscx4y3ihe3r6dk6ynsomrto6qir24pt45u7ihqgel5nlkpcf3xw"}})
	resources.Add("/MessageTestValueTransferSimplefailtotransfermorefundsthansenderbalance0", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(124), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacecf2zu6l2wcj63drvhkths2zg7hknmhcdmgr6hso7xq2ul6kjvvgs"}})
	resources.Add("/MessageTestValueTransferSimplefailtotransfermorefundsthansenderhaswhensenderbalancezero", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}, Penalty: abi.NewTokenAmount(120), Reward: abi.NewTokenAmount(0), Root: "bafy2bzacebfvultk6btvh6iggxrwk22zupu7w3nty2xgnfgkqux5fpvd2fae2"}})
	resources.Add("/MessageTestValueTransferSimplesuccessfullytransferfundsfromsendertoreceiver", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(130), Root: "bafy2bzacea4c3sc34mxpsmszmwxn6mkbanfoonl7kkpbrmkrwph54duqaajkq"}})
	resources.Add("/MessageTestValueTransferSimplesuccessfullytransferzerofundsfromsendertoreceiver", []types.ApplyMessageResult{types.ApplyMessageResult{Receipt: types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 126}, Penalty: abi.NewTokenAmount(0), Reward: abi.NewTokenAmount(126), Root: "bafy2bzaceceuzxvv3wqzxrw7kkj6asnexsyaxwh64nvtm5tptnimixbhfs63a"}})
	resources.Add("/TipSetTestBlockMessageApplicationSECPandBLSmessagescostdifferentamountsofgas", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 246}, types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 268}}, Root: "bafy2bzaceb6nh2igke3kzzbbhpxpma3h4xrtwirawict6nj5orhnmmbwh2nnu"}})
	resources.Add("/TipSetTestBlockMessageDeduplicationapplyaduplicatedBLSmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}}, Root: "bafy2bzacedjjmzcicn7o2jszb7djzzta5nkshzomfeou2ujpjcslhyzeljeb4"}})
	resources.Add("/TipSetTestBlockMessageDeduplicationapplyasingleBLSmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}}, Root: "bafy2bzacedjjmzcicn7o2jszb7djzzta5nkshzomfeou2ujpjcslhyzeljeb4"}})
	resources.Add("/TipSetTestBlockMessageDeduplicationapplyasingleSECPmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 268}}, Root: "bafy2bzaceajk2c3boonghz3vmiiit66linpgeeiieorc4t2wsrnho5whadchm"}})
	resources.Add("/TipSetTestBlockMessageDeduplicationapplyduplicateBLSandSECPmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 54}}, Root: "bafy2bzacebeba6zho3xqip5dgad5yvhhiraxkjpjrcx47ndujx26rnobeagrm"}})
	resources.Add("/TipSetTestBlockMessageDeduplicationapplyduplicateSECPmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 268}}, Root: "bafy2bzaceajk2c3boonghz3vmiiit66linpgeeiieorc4t2wsrnho5whadchm"}})
	resources.Add("/TipSetTestMinerRewardsAndPenaltiesinsufficientgastocoverreturnvalue", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{0x82, 0x42, 0x0, 0x64, 0x42, 0x0, 0x65}, GasUsed: 522}}, Root: "bafy2bzaceap7ay3vpfc3wuevmbrtwzhztq2gd6hpjsckp2jsq2ijn46v7w736"}, types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 7, ReturnValue: []uint8{}, GasUsed: 521}}, Root: "bafy2bzacedpp5wcyjud2mppyddtqzhxuqgh7jrpinceip5dfhxmuordk2fg72"}})
	resources.Add("/TipSetTestMinerRewardsAndPenaltiesminerpenaltyexceedsdeclaredgaslimitforBLSmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzaceczzux2kgicwa7iftuxpxoiv3jfz5zmpw756xdaz74xrt3nysyz7i"}})
	resources.Add("/TipSetTestMinerRewardsAndPenaltiesminerpenaltyexceedsdeclaredgaslimitforSECPmessage", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzaceaiu3lud7lrb55jnaoe33cjq5ue2wy6rm3wf4nzvvliuww6vucbty"}})
	resources.Add("/TipSetTestMinerRewardsAndPenaltiesoksimplesend", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}, types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 130}}, Root: "bafy2bzacecwrang5j7jr45tifamuznxtasw72lwiqejyoysuocqieqjzw7ede"}, types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 92}, types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 92}}, Root: "bafy2bzacecxcnidli4dkeywtgbmxnolyv2ugvip6tp33c6epzr6gkfsqddnxa"}, types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 92}, types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 92}}, Root: "bafy2bzaceduftevycbnrxvul2wdijadvqo6kiapyqomqhoyxcw7jvude7cnti"}, types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 54}, types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 54}}, Root: "bafy2bzaceblxcvp5nywbn6pomd6entrtbqvu5nqdqr7v2yrj5amstodowk6ty"}})
	resources.Add("/TipSetTestMinerRewardsAndPenaltiespenalizesenderdoestexist", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzaceb4ebhrwp65j2wmcvhssqnddbikqsc5m7rucxsrlusbt4nmv5ov6g"}})
	resources.Add("/TipSetTestMinerRewardsAndPenaltiespenalizesenderinsufficientbalance", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 0, ReturnValue: []uint8{}, GasUsed: 58}, types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzacecawiyijyqw23lq5nqgrih4ghmxa5oa2kgdii6hregn3kwcjsmm6u"}})
	resources.Add("/TipSetTestMinerRewardsAndPenaltiespenalizesendernonaccount", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}, types.MessageReceipt{ExitCode: 1, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzacebukzstfz5yssrrhmki5i5h27tu5q22r5oebhzaldngnodavcjxcs"}})
	resources.Add("/TipSetTestMinerRewardsAndPenaltiespenalizewrongcallseqnum", []types.ApplyTipSetResult{types.ApplyTipSetResult{Receipts: []types.MessageReceipt{types.MessageReceipt{ExitCode: 2, ReturnValue: []uint8{}, GasUsed: 0}}, Root: "bafy2bzaced4sw25rojgkqob747wejfqqeegw4sx6eqlj5alp66womjf6elr6u"}})
}
