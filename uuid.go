package uuid

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUID pgx type wrapper for google/uuid.UUID
type UUID uuid.UUID

// ScanUUID implements pgtype.UUIDScanner interface
func (u *UUID) ScanUUID(v pgtype.UUID) error {
	if !v.Valid {
		return fmt.Errorf("cannot scan NULL into *uuid.UUID")
	}

	*u = v.Bytes
	return nil
}

// UUIDValue implements pgtype.UUIDValuer interface
func (u UUID) UUIDValue() (pgtype.UUID, error) {
	return pgtype.UUID{Bytes: u, Valid: true}, nil
}

// NullUUID pgx type wrapper for google/uuid.NullUUID
type NullUUID uuid.NullUUID

// ScanUUID implements pgtype.UUIDScanner interface
func (u *NullUUID) ScanUUID(v pgtype.UUID) error {
	*u = NullUUID{UUID: v.Bytes, Valid: v.Valid}
	return nil
}

// UUIDValue implements pgtype.UUIDValuer interface
func (u NullUUID) UUIDValue() (pgtype.UUID, error) {
	return pgtype.UUID{Bytes: u.UUID, Valid: u.Valid}, nil
}

// TryWrapUUIDEncodePlan implements pgtype.TryWrapEncodePlanFunc interface
func TryWrapUUIDEncodePlan(value any) (plan pgtype.WrappedEncodePlanNextSetter, nextValue any, ok bool) {
	switch value := value.(type) {
	case uuid.UUID:
		return &wrapUUIDEncodePlan{}, UUID(value), true
	case uuid.NullUUID:
		return &wrapNullUUIDEncodePlan{}, NullUUID(value), true
	}

	return nil, nil, false
}

type wrapUUIDEncodePlan struct {
	next pgtype.EncodePlan
}

// SetNext implements pgtype.WrappedEncodePlanNextSetter interface
func (plan *wrapUUIDEncodePlan) SetNext(next pgtype.EncodePlan) { plan.next = next }

// Encode implements pgtype.EncodePlan interface
func (plan *wrapUUIDEncodePlan) Encode(value any, buf []byte) (newBuf []byte, err error) {
	return plan.next.Encode(UUID(value.(uuid.UUID)), buf)
}

type wrapNullUUIDEncodePlan struct {
	next pgtype.EncodePlan
}

// SetNext implements pgtype.WrappedEncodePlanNextSetter interface
func (plan *wrapNullUUIDEncodePlan) SetNext(next pgtype.EncodePlan) { plan.next = next }

// Encode implements pgtype.EncodePlan interface
func (plan *wrapNullUUIDEncodePlan) Encode(value any, buf []byte) (newBuf []byte, err error) {
	return plan.next.Encode(NullUUID(value.(uuid.NullUUID)), buf)
}

// TryWrapUUIDScanPlan implements pgtype.TryWrapScanPlanFunc
func TryWrapUUIDScanPlan(target any) (plan pgtype.WrappedScanPlanNextSetter, nextDst any, ok bool) {
	switch target := target.(type) {
	case *uuid.UUID:
		return &wrapUUIDScanPlan{}, (*UUID)(target), true
	case *uuid.NullUUID:
		return &wrapNullUUIDScanPlan{}, (*NullUUID)(target), true
	}

	return nil, nil, false
}

type wrapUUIDScanPlan struct {
	next pgtype.ScanPlan
}

// SetNext implements pgtype.WrappedScanPlanNextSetter interface
func (plan *wrapUUIDScanPlan) SetNext(next pgtype.ScanPlan) { plan.next = next }

// Scan implements pgtype.ScanPlan interface
func (plan *wrapUUIDScanPlan) Scan(src []byte, dst any) error {
	return plan.next.Scan(src, (*UUID)(dst.(*uuid.UUID)))
}

type wrapNullUUIDScanPlan struct {
	next pgtype.ScanPlan
}

// SetNext implements pgtype.WrappedScanPlanNextSetter interface
func (plan *wrapNullUUIDScanPlan) SetNext(next pgtype.ScanPlan) { plan.next = next }

// Scan implements pgtype.ScanPlan interface
func (plan *wrapNullUUIDScanPlan) Scan(src []byte, dst any) error {
	return plan.next.Scan(src, (*NullUUID)(dst.(*uuid.NullUUID)))
}

// UUIDCodec pgx type wrapper for pgtype.Codec
//revive:disable-next-line:exported
type UUIDCodec struct {
	pgtype.UUIDCodec
}

// DecodeValue implements pgtype.Codec interface
func (UUIDCodec) DecodeValue(tm *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}

	var target uuid.UUID
	scanPlan := tm.PlanScan(oid, format, &target)
	if scanPlan == nil {
		return nil, fmt.Errorf("PlanScan did not find a plan")
	}

	err := scanPlan.Scan(src, &target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

// Register registers the github.com/google/uuid integration with a pgtype.Map.
func Register(tm *pgtype.Map) {
	tm.TryWrapEncodePlanFuncs = append([]pgtype.TryWrapEncodePlanFunc{TryWrapUUIDEncodePlan}, tm.TryWrapEncodePlanFuncs...)
	tm.TryWrapScanPlanFuncs = append([]pgtype.TryWrapScanPlanFunc{TryWrapUUIDScanPlan}, tm.TryWrapScanPlanFuncs...)

	tm.RegisterType(&pgtype.Type{
		Name:  "uuid",
		OID:   pgtype.UUIDOID,
		Codec: UUIDCodec{},
	})
}
