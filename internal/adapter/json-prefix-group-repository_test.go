package adapter

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"lbmaster-advanced-groups-api/internal/domain"
	"os"
)

var _ = Describe("PrefixRepositoryTest", func() {
	var repo domain.PrefixGroupRepository
	var testPath string

	BeforeEach(func() {
		testPath = fmt.Sprintf("../../test/%s.json", uuid.New().String())
		c, err := ioutil.ReadFile("../../test/example.config.json")
		Expect(err).ToNot(HaveOccurred())
		err = ioutil.WriteFile(testPath, c, 0644)
		Expect(err).ToNot(HaveOccurred())
		repo = NewJsonPrefixGroupRepository(testPath)
	})

	AfterEach(func() {
		Expect(os.Remove(testPath)).ToNot(HaveOccurred())
	})

	It("returns a list of prefix groups", func() {
		p, err := repo.List()

		Expect(err).ToNot(HaveOccurred())
		Expect(p).To(HaveLen(1))
		Expect(p[0].Index).To(Equal(0))
		Expect(p[0].Prefix).To(Equal("[VIP] "))
	})

	It("returns not found for not existing prefix group", func() {
		_, err := repo.Members(domain.PrefixGroup{
			Index:  2,
			Prefix: "somePrefix",
		})

		Expect(err).To(MatchError(domain.ErrPrefixGroupNotFound))
	})

	It("lists members of a prefix group", func() {
		m, err := repo.Members(domain.PrefixGroup{
			Index:  0,
			Prefix: "somePrefix",
		})

		Expect(err).ToNot(HaveOccurred())
		Expect(m).To(HaveLen(2))
		Expect(m[0]).To(Equal(domain.SteamUID("76561111111111111")))
		Expect(m[1]).To(Equal(domain.SteamUID("76561122222222222")))
	})

	It("adds a new member to a prefix group", func() {
		group := domain.PrefixGroup{
			Index:  0,
			Prefix: "somePrefix",
		}
		err := repo.AddMember(group, "76561133333333333")
		Expect(err).ToNot(HaveOccurred())

		m, err := repo.Members(group)
		Expect(err).ToNot(HaveOccurred())
		Expect(m).To(HaveLen(3))
		Expect(m[0]).To(Equal(domain.SteamUID("76561111111111111")))
		Expect(m[1]).To(Equal(domain.SteamUID("76561122222222222")))
		Expect(m[2]).To(Equal(domain.SteamUID("76561133333333333")))
	})

	It("does not add member if they already exist", func() {
		group := domain.PrefixGroup{
			Index:  0,
			Prefix: "somePrefix",
		}
		err := repo.AddMember(group, "76561122222222222")
		Expect(err).ToNot(HaveOccurred())

		m, err := repo.Members(group)
		Expect(err).ToNot(HaveOccurred())
		Expect(m).To(HaveLen(2))
	})

	It("removes a member from a prefix group", func() {
		group := domain.PrefixGroup{
			Index:  0,
			Prefix: "somePrefix",
		}
		err := repo.RemoveMember(group, "76561122222222222")
		Expect(err).ToNot(HaveOccurred())

		m, err := repo.Members(group)
		Expect(err).ToNot(HaveOccurred())
		Expect(m).To(HaveLen(1))
		Expect(m[0]).To(Equal(domain.SteamUID("76561111111111111")))
	})
})
