// Copyright 2015 The MOAC-core Authors
// This file is part of the MOAC-core library.
//
// The MOAC-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The MOAC-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the MOAC-core library. If not, see <http://www.gnu.org/licenses/>.

package params

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main MoacNode network.

var MainnetBootnodes = []string{
	// MoacNode Foundation Go Bootnodes
	//BeijingEnode, Pangu boot nodes:
	// "enode://25fb17611768ed9a7ed7dc9815f75316478162efe5c99df67b5b5e9b7343de114a9baaf3c5a5bc7ae7459df9b06869646eb6e212256c239f7dc29fd5d9a26352@139.198.16.48:30333",
	// //Shanghai node
	// "enode://b81db3b00b8aab3b1e33381e23ade1bafb6d41a1c084df14240a93fe98ab2043e58d0ccac43c850ea912c1231a939be5421b1d34c03f5add868cfdb2adbaa359@139.198.177.169:30333",
	// //AMAZON nodes
	// "enode://56d1df4beff7834fa0f91ffaf013d7c0851ce9aa8453f26b148e30941f4a1a9ba43e1e125caf6cacf1f0ad2571c1e722c9b64456240a73af776195f4c3a8c72e@18.217.180.94:30333",
	// "enode://e3a5b39ad5019543fd798deb911c65882325800d60ef288aef1822a047172222bd3412ff0e5896259e14af34c6a8520d056b0732e23f22d6bcd5c31ed8585aa5@18.217.97.151:30333",
	// "enode://3f4eb607c1105b842dcd8c0e00f50f61beb8cd36f304e53c325dcc37afd65cb9e7c4ac075cd9d67895d7e5cfa875e3b4d8dc0f15f09b47135ceefeda7db5ca02@52.15.143.41:30333",
	//Nuwa bootnodes:
	"enode://3d8ba7cef2dcc8e25bd508e6a42d7b36a957afe8dc5cd9603ffd0857aae68b2768a06de6054816edec407a4579b15b5c7494e9c7e8323b015d6e9d518576e9eb@18.233.50.84:30333",
	"enode://6f29307c8715502c11791704759d88910e091f5d4a6daf509fc9d0d8653967fd2cda76c94731e2db6654ff889a684a55f9a6bc6f0b3144282eb2c3e7cbccb4b1@18.211.187.153:30333",
	"enode://1cf0eb6052ee692e85cf5447c9b1ef95c57f23e2efe1cd8cd5a01b44feadc3395f3aa02e0b1130dffcfdac3225b3ce8d773bfd148efadc3846353337a45c99a2@54.71.244.228:30333",
	"enode://bebab9587e939fab7e21922864cdf6668193e1fcfe5f69ef128216c96e5e08dae1a8d9b99e0cb158bd040c9ac249c3569d43675f0e637fb2cf9e2df9e39738d7@52.34.175.72:30333",
	//london, UK
	"enode://68c48e272a78e288c9e86ca2c272cf51ed9dbbe305e8f257cf347b4282dc04bc404ca25178025661216b4f15d610b1d536ab3c8610ca4357251ba85451a889c4@18.130.240.247:30333",
	//sydney, AU
	"enode://624c26083a9be2ac0624bc5e4d88633aca89f62bbca1e7929a05bfc9f06a42202168b12ed5303f95e26296f8092d07b61071eedaccfe7027ceb0e7330589bc5c@13.211.142.153:30333",
	//Canada central
	"enode://4c8fac5ed1f8e5676d8f6b722fdb69b289dd58fe5844f0ae036b9d257bc6c4ebb94641cd590f8d7ed282c7690217eaf17898712b974ce23cfccbe167edc1faa0@35.182.237.1:30333",
	//Singapore
	"enode://ab485ba14e3955cedeff78ef760bc34c450f80396ea5f2c9a2669e2aa7d16dec014a6772a5eacd3726f8ee83cc0877ce2c0cf2a9a95b6e3e239df15bcbe52d58@47.88.237.205:30333",
	//Tokyo, Japan
	"enode://3303754238a2a2f07aadf7683224768c562822517e270b38fad7cf6df9c1af6d94a42ce361bbe965a38f2a2b1dfc2873518a8c20a5a004e8aa105d3bc3dd4a8b@47.74.9.125:30333",
	// Guangdong, China
	"enode://9675508dfc11685d762f0b95f9f939af872baf2d2f289e6b142768030dd62a8f413fc7baa6bd2744a4ec7bf8d001ed59dc82e47e46702092097dcc7a7f08beee@139.198.122.215:30333",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Pangu 0.8.1 test network.
var TestnetBootnodes = []string{
	//pangu 0.8.2 testnetid = 101
	//md02:"enode://1a8e0758b99f5b1eb248b240ddc8efd6f4b65821ac96ac19dd4287fd8fbb23db84e257b4fa5a9473c0a21eef9a70db7504a5e7a026e6b34398ce4ff856e4ca0b@35.185.8.110:30333"
	//48f08596c760ee663c0ac3b48b5bf7476771b87586df7ea424f5bd3b6e2c40bc2089b5ebaaf364b6cf86e446f027ee083d8ddbd019e355877d8a1fae173d8a5b@
	//NUWA testnet bootnodes
	//China Guangzhou
	"enode://e88a288b19e55a79cdeabd9c5326602dbb2ba5a5e29eca66bdddf20d598cf7126fce9080b53399ef02e200891cd590b6ff9c23a8c15789f04790d1e69054bd9c@39.108.79.40:30333",
	//China Shanghai
	"enode://fcb843289fb0b692fc9d2449e8dce7b9143b735263f99c594d8c5c196f1cff9066812b5b978cca13edf17efa712f0e3e776f6880f92bc5fd58d0d27e3239f8d6@47.98.255.26:30333",
	//AMAZON AWS,
	//t6
	"enode://154ce8ba69e7a24c237a4c69928c814b4341877e623d9a0ce373ce97132d098c3c671122049ae85e5d081e8763a4d0f38178736457bd7197a6b75a49fa50156d@107.155.107.194:30333",
	//t4
	"enode://a44e287b1a3e2f6e898d92469242b42360f05f701594bc098b21930799ef41f92fe6c51cc51a1ca49ff4b887c63abd9c7a9786fba4317ab851d5c6043829b3ac@107.155.125.100:30333",
	//t5 node
	"enode://77fa8dd1bd80002f16abb9ea9fe7073025ba20b3d494eedb0980a859200a45b2de585be9071cbeada2f49904ac888b6423e1c4bcde3e57e543c423a9fa68c5c3@107.155.125.100:30333",
}
