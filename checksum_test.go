package nmea_test

import (
	. "github.com/munnik/go-nmea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Checksum", func() {
	var (
		sentence Sentence
		err      error
		raw      string
	)
	Describe("Parsing without options", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw)
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*21"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("returns not nil", func() {
				Expect(sentence).NotTo(BeNil())
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*25"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [21 != 25]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
		Context("a sentence with an empty checksum", func() {
			BeforeEach(func() {
				raw = "$23RSA,-028.8,A,-028.8,A*"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [41 != ]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Parsing with AllowEmptyChecksum", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw, "AllowEmptyChecksum")
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*21"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("returns not nil", func() {
				Expect(sentence).NotTo(BeNil())
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*25"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [21 != 25]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
		Context("a sentence with an empty checksum", func() {
			BeforeEach(func() {
				raw = "$23RSA,-028.8,A,-028.8,A*"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("returns not nil", func() {
				Expect(sentence).NotTo(BeNil())
			})
		})
	})
	Describe("Parsing with AllowChecksumMismatch", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw, "AllowChecksumMismatch")
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*21"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("returns not nil", func() {
				Expect(sentence).NotTo(BeNil())
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*25"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("returns not nil", func() {
				Expect(sentence).NotTo(BeNil())
			})
		})
		Context("a sentence with an empty checksum", func() {
			BeforeEach(func() {
				raw = "$23RSA,-028.8,A,-028.8,A*"
			})
			It("returns an error", func() {
				Expect(err).To(MatchError("nmea: sentence checksum mismatch [41 != ]"))
			})
			It("returns nil", func() {
				Expect(sentence).To(BeNil())
			})
		})
	})
	Describe("Parsing with AllowEmptyChecksum and AllowChecksumMismatch", func() {
		JustBeforeEach(func() {
			sentence, err = Parse(raw, "AllowEmptyChecksum", "AllowChecksumMismatch")
		})
		Context("a valid sentence", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*21"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("returns not nil", func() {
				Expect(sentence).NotTo(BeNil())
			})
		})
		Context("a sentence with a bad checksum", func() {
			BeforeEach(func() {
				raw = "$23DBS,01.9,f,0.58,M,00.3,F*25"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("returns not nil", func() {
				Expect(sentence).NotTo(BeNil())
			})
		})
		Context("a sentence with an empty checksum", func() {
			BeforeEach(func() {
				raw = "$23RSA,-028.8,A,-028.8,A*"
			})
			It("returns no errors", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("returns not nil", func() {
				Expect(sentence).NotTo(BeNil())
			})
		})
	})
})
