package repositories

import (
	"crypto/md5"
	//	"strconv"
	"time"

	"github.com/asdine/storm"
	"github.com/graph-uk/cafe-runner/data/models"
)

type Testpacks struct {
	Tx storm.Node
}

func (t *Testpacks) Create(zipfile *[]byte) *models.Testpack {
	hasher := md5.New()
	_, err := hasher.Write(*zipfile)
	check(err)

	testpack := &models.Testpack{
		Status:     models.TPStatusReadyForUnzip,
		UploadTime: time.Now(),
		Zip:        *zipfile,
		ZipHash:    hasher.Sum(nil),
	}

	check(t.Tx.Save(testpack))

	//docker run --rm -v D:\GraphWorking\cafeRunner\src\github.com\graph-uk\cafe-runner\tests\integration\static\testcafe:/tests -it --workdir /tests --entrypoint /bin/sh testcafe/testcafe -c "npm install"

	return testpack
}

func (t *Testpacks) Find(id int) *models.Testpack {
	res := &models.Testpack{}
	check(t.Tx.One(`ID`, id, res))
	return res
}

func (t *Testpacks) FindAll() *[]models.Testpack {
	res := &[]models.Testpack{}
	check(t.Tx.All(res))
	return res
}

func (t *Testpacks) Update(tespack *models.Testpack) {
	check(t.Tx.Update(tespack))
}
