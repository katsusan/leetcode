package godemo

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func bcryptCheck() {
	hp, err := bcrypt.GenerateFromPassword([]byte("mypwd"), 0)
	if err != nil {
		log.Fatalln("failed to generate bcrypt string")
	}
	log.Println("succeed to generate crpyted string:", string(hp))

	hp2, err := bcrypt.GenerateFromPassword([]byte("mypwd"), 0)
	if err != nil {
		log.Fatalln("failed to generate bcrypt string for 'mypwd'")
	}
	log.Println("succeed to generate crpyted string:", string(hp2))

	if err := bcrypt.CompareHashAndPassword(hp, []byte("mypassword")); err != nil {
		log.Println("not match, test OK")
	}

	if err := bcrypt.CompareHashAndPassword(hp, hp2); err == nil {
		log.Println("password match, test OK")
	}

	if err := bcrypt.CompareHashAndPassword(hp, []byte("hahha")); err == nil {
		log.Println("this log shouldn't be printed, because correct password is 'mypwd'")
	}
}
