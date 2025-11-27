package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mariotiara/sfe-data-pipe/configs"
	"github.com/mariotiara/sfe-data-pipe/internal/domain/ezengagecalldetail"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
)

type EZEngageCallDetailRepository struct {
	logger logger.Logger
	db     *sql.DB
	config *configs.Config
}

func NewEZEngageCallDetailRepository(db *sql.DB, config *configs.Config, logger logger.Logger) *EZEngageCallDetailRepository {
	return &EZEngageCallDetailRepository{db: db, config: config, logger: logger}
}

func (r *EZEngageCallDetailRepository) RemoveThisMonthData(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		DELETE FROM ez_engage_call_detail
		WHERE EXTRACT (YEAR FROM created_at)=$1
		AND EXTRACT (MONTH FROM created_at)=$2
		RETURNING 1
	`

	rows, err := r.db.QueryContext(ctx, query, year, month)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	r.logger.Info(ctx, fmt.Sprintf("%d Rows has been removed", count))
	return count, nil
}

func (r *EZEngageCallDetailRepository) HasThisMonthData(ctx context.Context) (bool, error) {
	now := time.Now().UTC()
	year, month := now.Year(), now.Month()

	query := `
		SELECT COUNT(1)
		FROM ez_engage_call_detail
		WHERE EXTRACT (YEAR FROM created_at)=$1
		AND EXTRACT (MONTH FROM created_at)=$2
	`
	var count int

	err := r.db.QueryRowContext(ctx, query, year, month).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func (r *EZEngageCallDetailRepository) Save(ctx context.Context, e *ezengagecalldetail.EZEngageCallDetail) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Using 76 columns
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO ez_engage_call_detail (
			date TEXT,
			team TEXT,
			position TEXT,
			business_type TEXT,
			sales_org TEXT,
			territory_code TEXT,
			sales_representative TEXT,
			status TEXT,
			missed_call_remark TEXT,
			planned_unplanned TEXT,
			activity TEXT,
			otherwork_name TEXT,
			otherwork_note TEXT,
			customer_code TEXT,
			customer_name TEXT,
			customer_ship_to_code TEXT,
			customer_ship_to_name TEXT,
			customer_city TEXT,
			customer_class TEXT,
			frequency TEXT,
			principal TEXT,
			product TEXT,
			product_code TEXT,
			stock_inventory_quantity TEXT,
			listed TEXT,
			shelf_stock TEXT,
			out_of_stock TEXT,
			location TEXT,
			shelf_space TEXT,
			share_of_space TEXT,
			no_of_facings TEXT,
			no_of_photos_upload TEXT,
			product_enlistment TEXT,
			promotional_activity TEXT,
			sales_recommendation_from_ico TEXT,
			person_in_charge TEXT,
			with_contract TEXT,
			start_date TEXT,
			end_date TEXT,
			virtual_images TEXT,
			presentation_duration_time TEXT,
			presentation_file_name TEXT,
			perform_collection TEXT,
			mode_of_collection TEXT,
			placed_order TEXT,
			mode_of_order TEXT,
			reason_of_not_using_ezrx TEXT,
			check_in_time TEXT,
			check_in_date TEXT,
			check_out_time TEXT,
			check_out_date TEXT,
			visit_duration TEXT,
			travel_duration_time TEXT,
			pre_call_notes TEXT,
			post_call_notes TEXT,
			barrier_encountered TEXT,
			call_source TEXT,
			route_status TEXT,
			modality_of_call TEXT,
			signature_image TEXT,
			signature_count TEXT,
			signature_start_time TEXT,
			signature_end_time TEXT,
			signature_diff_time TEXT,
			longitude_data TEXT,
			latitude_data TEXT,
			actual_longitude_data TEXT,
			actual_latitude_data TEXT,
			location_accuracy TEXT,
			geo_location_remarks TEXT,
			mode TEXT,
			work_with TEX
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,
			$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,
			$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,
			$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,
			$47,$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,
			$58,$59,$60,$61,$62,$63,$64,$65,$66,$67,$68,
			$69,$70,$71,$72
		)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		e.Date, e.Team, e.Position, e.BusinessType, e.SalesOrg, e.TerritoryCode, e.SalesRepresentative, e.Status,
		e.MissedCallRemark, e.PlannedUnplanned, e.Activity, e.OtherworkName, e.OtherworkNote,
		e.CustomerCode, e.CustomerName, e.CustomerShipToCode, e.CustomerShipToName, e.CustomerCity,
		e.CustomerClass, e.Frequency, e.Principal, e.Product, e.ProductCode, e.StockInventoryQuantity,
		e.Listed, e.ShelfStock, e.OutOfStock, e.Location, e.ShelfSpace, e.ShareOfSpace, e.NoOfFacings,
		e.NoOfPhotosUpload, e.ProductEnlistment, e.PromotionalActivity, e.SalesRecommendationFromICO,
		e.PersonInCharge, e.WithContract, e.StartDate, e.EndDate, e.VirtualImages, e.PresentationDurationTime,
		e.PresentationFileName, e.PerformCollection, e.ModeOfCollection, e.PlacedOrder, e.ModeOfOrder,
		e.ReasonOfNotUsingEzrx, e.CheckInTime, e.CheckInDate, e.CheckOutTime, e.CheckOutDate, e.VisitDuration,
		e.TravelDurationTime, e.PreCallNotes, e.PostCallNotes, e.BarrierEncountered, e.CallSource, e.RouteStatus,
		e.ModalityOfCall, e.SignatureImage, e.SignatureCount, e.SignatureStartTime, e.SignatureEndTime,
		e.SignatureDiffTime, e.LongitudeData, e.LatitudeData, e.ActualLongitudeData, e.ActualLatitudeData,
		e.LocationAccuracy, e.GeoLocationRemarks, e.Mode, e.WorkWith,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *EZEngageCallDetailRepository) SaveRange(ctx context.Context, details []*ezengagecalldetail.EZEngageCallDetail) error {
	batchSize := r.config.DBBatchSize
	tinserted := 0
	for i := 0; i < len(details); i += batchSize {
		end := i + batchSize
		if end > len(details) {
			end = len(details)
		}
		batch := details[i:end]

		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}

		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO ez_engage_call_detail (
				date, team, position, business_type, sales_org, territory_code, sales_representative, status,
				missed_call_remark, planned_unplanned, activity, otherwork_name, otherwork_note,
				customer_code, customer_name, customer_ship_to_code, customer_ship_to_name, customer_city,
				customer_class, frequency, principal, product, product_code, stock_inventory_quantity,
				listed, shelf_stock, out_of_stock, location, shelf_space, share_of_space, no_of_facings,
				no_of_photos_upload, product_enlistment, promotional_activity, sales_recommendation_from_ico,
				person_in_charge, with_contract, start_date, end_date, virtual_images, presentation_duration_time,
				presentation_file_name, perform_collection, mode_of_collection, placed_order, mode_of_order,
				reason_of_not_using_ezrx, check_in_time, check_in_date, check_out_time, check_out_date, visit_duration,
				travel_duration_time, pre_call_notes, post_call_notes, barrier_encountered, call_source, route_status,
				modality_of_call, signature_image, signature_count, signature_start_time, signature_end_time,
				signature_diff_time, longitude_data, latitude_data, actual_longitude_data, actual_latitude_data,
				location_accuracy, geo_location_remarks, mode, work_with
			) VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
				$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,
				$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,$58,$59,$60,
				$61,$62,$63,$64,$65,$66,$67,$68,$69,$70,$71,$72
			)
		`)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("prepare stmt: %w", err)
		}

		for _, e := range batch {
			_, err := stmt.ExecContext(ctx,
				e.Date, e.Team, e.Position, e.BusinessType, e.SalesOrg, e.TerritoryCode, e.SalesRepresentative, e.Status,
				e.MissedCallRemark, e.PlannedUnplanned, e.Activity, e.OtherworkName, e.OtherworkNote,
				e.CustomerCode, e.CustomerName, e.CustomerShipToCode, e.CustomerShipToName, e.CustomerCity,
				e.CustomerClass, e.Frequency, e.Principal, e.Product, e.ProductCode, e.StockInventoryQuantity,
				e.Listed, e.ShelfStock, e.OutOfStock, e.Location, e.ShelfSpace, e.ShareOfSpace, e.NoOfFacings,
				e.NoOfPhotosUpload, e.ProductEnlistment, e.PromotionalActivity, e.SalesRecommendationFromICO,
				e.PersonInCharge, e.WithContract, e.StartDate, e.EndDate, e.VirtualImages, e.PresentationDurationTime,
				e.PresentationFileName, e.PerformCollection, e.ModeOfCollection, e.PlacedOrder, e.ModeOfOrder,
				e.ReasonOfNotUsingEzrx, e.CheckInTime, e.CheckInDate, e.CheckOutTime, e.CheckOutDate, e.VisitDuration,
				e.TravelDurationTime, e.PreCallNotes, e.PostCallNotes, e.BarrierEncountered, e.CallSource, e.RouteStatus,
				e.ModalityOfCall, e.SignatureImage, e.SignatureCount, e.SignatureStartTime, e.SignatureEndTime,
				e.SignatureDiffTime, e.LongitudeData, e.LatitudeData, e.ActualLongitudeData, e.ActualLatitudeData,
				e.LocationAccuracy, e.GeoLocationRemarks, e.Mode, e.WorkWith,
			)
			if err != nil {
				stmt.Close()
				tx.Rollback()
				return fmt.Errorf("exec batch insert: %w", err)
			}
		}

		stmt.Close()
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit batch: %w", err)
		}

		tinserted = end
	}

	r.logger.Info(ctx, fmt.Sprintf("%d of %d rows successfully inserted", len(details), tinserted))
	return nil
}
