package memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/ellcrys/util"
	"github.com/ellcrys/cocoon/core/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMemoryLock(t *testing.T) {
	Convey("MemoryLock", t, func() {
		Convey(".AcquireLock", func() {
			Convey("Should successfully acquire a lock", func() {
				key := util.RandString(10)
				l := NewLock(key)
				err := l.Acquire()
				So(err, ShouldBeNil)

				Convey("Should have no problem re-acquiring a lock as long as TTL has not passed", func() {
					err := l.Acquire()
					So(err, ShouldBeNil)
				})

				Convey("Should fail if lock has already been acquired by a different session", func() {
					l := NewLock(key)
					err := l.Acquire()
					So(err, ShouldResemble, types.ErrLockAlreadyAcquired)
				})
			})
		})

		Convey(".ReleaseLock", func() {
			Convey("Should successfully release an acquired lock", func() {
				key := util.RandString(10)
				l := NewLock(key)
				err := l.Acquire()
				So(err, ShouldBeNil)
				err = l.Release()
				So(err, ShouldBeNil)

				Convey("Should successfully acquire a released lock", func() {
					err := l.Acquire()
					So(err, ShouldBeNil)
				})
			})
		})

		Convey(".IsAcquirer", func() {

			Convey("Should return error if lock has no previously acquired key", func() {
				l := NewLock("")
				err := l.IsAcquirer()
				So(err, ShouldResemble, fmt.Errorf("key is not set"))
			})

			Convey("Should return nil if lock is still the acquirer of a lock on it's key", func() {
				key := util.RandString(10)
				l := NewLock(key)
				err := l.Acquire()
				So(err, ShouldBeNil)
				err = l.IsAcquirer()
				So(err, ShouldBeNil)
			})

			Convey("Should return err if lock is no longer acquired due to TTL being reached", func() {
				key := util.RandString(10)
				l := NewLock(key)
				LockTTL = 2 * time.Second
				err := l.Acquire()
				So(err, ShouldBeNil)
				cancel := StartLockWatcher()
				time.Sleep(4 * time.Second)
				cancel()
				err = l.IsAcquirer()
				So(err, ShouldResemble, types.ErrLockNotAcquired)
			})

		})
	})
}
