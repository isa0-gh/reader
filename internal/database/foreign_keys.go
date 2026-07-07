package database

import "gorm.io/gorm"

// EnsureForeignKeys adds FK constraints that are intentionally kept out of
// GORM's own AutoMigrate-driven relationship graph, either because gorm
// doesn't let an existing constraint be altered in place (series' cover
// image needs ON DELETE SET NULL, but a constraint created without one
// can't be changed by re-running AutoMigrate) or because declaring the
// relation to gorm at all would reintroduce a migration-order cycle
// (chapters -> series -> s3_objects -> chapters, since Series.CoverImage
// depends on S3Object and S3Object.Pages depends on Chapter). Running
// these as plain ALTER TABLE statements after every model has already been
// migrated sidesteps both problems. Safe to call on every boot: each
// constraint is dropped if present and re-added, so it always converges on
// the desired definition regardless of starting state.
func EnsureForeignKeys(db *gorm.DB) error {
	constraints := []struct {
		table string
		name  string
		add   string
	}{
		{
			table: "series",
			name:  "fk_series_cover_image",
			add:   "ALTER TABLE series ADD CONSTRAINT fk_series_cover_image FOREIGN KEY (cover_image_id) REFERENCES s3_objects(id) ON DELETE SET NULL",
		},
		{
			table: "chapters",
			name:  "fk_chapters_series",
			add:   "ALTER TABLE chapters ADD CONSTRAINT fk_chapters_series FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE",
		},
	}

	for _, c := range constraints {
		if err := db.Exec("ALTER TABLE " + c.table + " DROP CONSTRAINT IF EXISTS " + c.name).Error; err != nil {
			return err
		}
		if err := db.Exec(c.add).Error; err != nil {
			return err
		}
	}
	return nil
}
