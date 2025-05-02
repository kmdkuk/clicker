package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InputHandler", func() {
	var (
		handler *DefaultHandler
	)

	BeforeEach(func() {
		handler = &DefaultHandler{
			pressedKey: ebiten.KeyMeta, // Initialize with a default key
		}
	})

	Describe("NewDefaultKeyHandler", func() {
		It("should initialize with no pressed key", func() {
			Expect(handler.GetPressedKey()).To(Equal(KeyTypeNone))
		})
	})

	Describe("Update", func() {
		It("should not panic when updating", func() {
			Expect(func() {
				handler.Update()
			}).NotTo(Panic())
		})
	})

	Describe("getPressedKeyType", func() {
		It("should return the correct key type for UP", func() {
			ups := []ebiten.Key{
				ebiten.KeyArrowUp,
				ebiten.KeyW,
				ebiten.KeyK,
			}
			for _, up := range ups {
				handler.pressedKey = up
				keyType := handler.GetPressedKey()
				Expect(keyType).To(Equal(KeyTypeUp))
			}
		})

		It("should return the correct key type for DOWN", func() {
			downs := []ebiten.Key{
				ebiten.KeyArrowDown,
				ebiten.KeyS,
				ebiten.KeyJ,
			}
			for _, down := range downs {
				handler.pressedKey = down
				keyType := handler.GetPressedKey()
				Expect(keyType).To(Equal(KeyTypeDown))
			}
		})

		It("should return the correct key type for LEFT", func() {
			lefts := []ebiten.Key{
				ebiten.KeyArrowLeft,
				ebiten.KeyA,
				ebiten.KeyH,
			}
			for _, left := range lefts {
				handler.pressedKey = left
				keyType := handler.GetPressedKey()
				Expect(keyType).To(Equal(KeyTypeLeft))
			}
		})

		It("should return the correct key type for RIGHT", func() {
			rights := []ebiten.Key{
				ebiten.KeyArrowRight,
				ebiten.KeyD,
				ebiten.KeyL,
			}
			for _, right := range rights {
				handler.pressedKey = right
				keyType := handler.GetPressedKey()
				Expect(keyType).To(Equal(KeyTypeRight))
			}
		})

		It("should return the correct key type for Desicison", func() {
			desicions := []ebiten.Key{
				ebiten.KeyEnter,
				ebiten.KeySpace,
			}
			for _, desicion := range desicions {
				handler.pressedKey = desicion
				keyType := handler.GetPressedKey()
				Expect(keyType).To(Equal(KeyTypeDecision))
			}
		})

		It("should return NONE for other keys", func() {
			handler.pressedKey = ebiten.KeyMeta
			keyType := handler.GetPressedKey()
			Expect(keyType).To(Equal(KeyTypeNone))
		})
	})
})
