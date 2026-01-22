/*
Copyright 2021 Ethan Blanton <eblanton@buffalo.edu>

This file is part of a CSE 486/586 project from the University at
Buffalo.  Distribution of this file or its associated repository
requires the written permission of Ethan Blanton.  Sharing this file
may be a violation of academic integrity, please consult the course
policies for more package.
*/

package kdht

import "crypto/sha1"

// KeyBytes is the size in bytes of the keys in the k-DHT
const KeyBytes = sha1.Size

// KeyBits is the ssize in bits of the keys in the k-DHT
const KeyBits = 8 * KeyBytes
