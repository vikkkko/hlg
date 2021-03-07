package dao

import "github.com/tpkeeper/pitaya/models"

func (d *Dao) UpdateOrInsertHbt(mfi *models.Pitaya) error {
	return d.db.Save(mfi).Error
}

func (d *Dao) GetPitaya() (mfi *models.Pitaya, err error) {
	mfi = &models.Pitaya{}
	err = d.db.First(mfi).Error
	return
}
