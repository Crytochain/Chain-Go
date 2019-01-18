// Copyright 2017 The MOAC-core Authors
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

const (
	// These are the multipliers for MOAC denominations.
	// Example: To get the sha value of an amount in 'Grand', use
	//
	//    new(big.Int).Mul(value, big.NewInt(params.Grand))
	//
	Sha     = 1
	Femtomc = 1e3
	Picomc  = 1e6
	Xiao    = 1e9
	Sand    = 1e12
	Milli   = 1e15
	Mc      = 1e18
	Moac    = 1e18
	Grand   = 1e21
)
