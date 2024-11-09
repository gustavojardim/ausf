package akma

import (
	"fmt"

	"github.com/free5gc/util/ueauth"
)

/*
deriveAKMAKey is a local function that derives the AKMA key and A_KID from the provided K_AUSF key.
This function needs to follow the 3GPP specification.
*/
func DerivateAkmaKey(Kausf []byte, supi []byte, rid string, mcc string, mnc string) ([]byte, string, error) {
	// AKMA Key derivation as described in TS 33.535 Annex A.2
	P0_AKMA := []byte("AKMA")
	P1_AKMA := supi
	Kakma, err := ueauth.GetKDFValue(Kausf, "80", P0_AKMA, ueauth.KDFLen(P0_AKMA), P1_AKMA, ueauth.KDFLen(P1_AKMA))
	if err != nil {
		return nil, "", fmt.Errorf("AKMA Key generation failed: %+v", err)
	}

	// A-TID derivation as described in TS 33.535 Annex A.3
	P0_ATID := []byte("A-TID")
	P1_ATID := supi
	atid, err := ueauth.GetKDFValue(Kausf, "81", P0_ATID, ueauth.KDFLen(P0_ATID), P1_ATID, ueauth.KDFLen(P1_ATID))
	if err != nil {
		return nil, "", fmt.Errorf("A-TID derivation failed: %+v", err)
	}

	// A-KID defined in the TS 33.535 clause 6.1
	realm := mcc + "" + mnc
	atidStr := fmt.Sprintf("%x", atid)
	akId := atidStr + rid + "@" + realm

	return Kakma, akId, nil
}

/*
deriveAFKey derives an application-specific key (K_AF) from the AKMA key (K_AKMA) based on the AF_ID.
*/
func DerivateApplicationFunctionKey(Kakma []byte, afId string) ([]byte, error) {
	// Application Function Key derivation as described in TS 33.535 Annex A.4
	P0_Kaf := []byte(afId)
	Kaf, err := ueauth.GetKDFValue(Kakma, "82", P0_Kaf, ueauth.KDFLen(P0_Kaf))

	if err != nil {
		return nil, fmt.Errorf("Application Function Key generation failed: %+v", err)
	}

	return Kaf, nil
}
