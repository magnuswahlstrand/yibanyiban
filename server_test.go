package yibanyiban

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test the HTTP return codes for IBANHandler
func TestIBANHandler(t *testing.T) {
	tcs := []struct {
		name               string
		method             string
		path               string
		expectedStatusCode int
	}{
		{"ok #1", "GET", "validate?iban=GB82WEST12345698765432", http.StatusOK},
		{"ok #2", "GET", "validate?iban=AL85751639367318444714198669", http.StatusOK},
		{"invalid method", "POST", "validate?iban=AL85751639367318444714198669", http.StatusMethodNotAllowed},
		{"missing parameter iban", "GET", "validate?lolban=AL85751639367318444714198669", http.StatusForbidden},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "localhost:8080/"+tc.path, nil)
			rec := httptest.NewRecorder()

			ValidateIBANHandler(rec, req)
			res := rec.Result()

			if res.StatusCode != tc.expectedStatusCode {
				t.Fatalf("expected '%s', got '%s'\n", http.StatusText(tc.expectedStatusCode), http.StatusText(res.StatusCode))
			}
		})
	}
}

func TestValidateIBAN(t *testing.T) {
	tcs := []struct {
		IBAN     string
		expected bool
	}{
		{"GB82WEST12345698765432", true},
		{"AL86751639367318444714198669", true},
		{"AL87751639367318444714198669", false},
	}

	for _, tc := range tcs {
		valid, _ := validateIBAN(tc.IBAN)
		if valid != tc.expected {
			t.Errorf("expected validateIBAN(%q)=%t, got %t", tc.IBAN, tc.expected, valid)
		}
	}
}
func TestInvalidIBANError(t *testing.T) {
	tcs := []struct {
		IBAN        string
		expectedErr error
	}{
		{"GB8", errNumberTooShort},
		{"", errNumberTooShort},
		{"AL86751639367318444714198669AL86751639367318444714198669", errNumberTooLong},
		{"SE120312301023012301203012301230120301230210300123010230", errNumberTooLong},
		{"GB82WEST12345698765432&", errInvalidCharacters},
		{"GB82WEST12345698765432*", errInvalidCharacters},
		{"GB83WEST12345698765432", errCheckSumIncorrect},
		{"AL85751639367318444714198669", errCheckSumIncorrect},
		{"GB82WEST12345698765432", nil},
	}

	for _, tc := range tcs {
		_, err := validateIBAN(tc.IBAN)
		if err != tc.expectedErr {
			t.Errorf("expected error %q from validateIBAN(%q), got %q", tc.expectedErr, tc.IBAN, err)
		}
	}
}

// Validate against the list from https://www.mobilefish.com/download/iban/random_generated_iban.txt
func TestValidateManyIBAN(t *testing.T) {
	for _, IBAN := range validIBANs {
		valid, _ := validateIBAN(IBAN)
		if valid != true {
			t.Errorf("expected validateIBAN(%q)=%t, got %t", IBAN, true, valid)
		}
	}
}

var validIBANs = []string{"AL86751639367318444714198669",
	"AL89515635252277023782748302",
	"AL39153296222641598198140883",
	"AL47907501989147671525950076",
	"AL55398719849655505231753964",
	"AD2531377125214715353449",
	"AD9764782778017799549345",
	"AD4079739934060166934190",
	"AD3210446914824799260335",
	"AD1781438353588817727122",
	"AT582774098454337653",
	"AT220332087576467472",
	"AT328650112318219886",
	"AT193357281080332578",
	"AT535755326448639816",
	"BE16517682243567",
	"BE46943937718104",
	"BE75270187592710",
	"BE58465045170210",
	"BE49149522496291",
	"BA534130469841865537",
	"BA515388988295860588",
	"BA182655808222815318",
	"BA105531662061034080",
	"BA198940842595891985",
	"BG08NXYF73308507056085",
	"BG22OOPG05631394112384",
	"BG30XCMJ43923350257238",
	"BG66ZKSV96204746173581",
	"BG62TOZJ59790808155256",
	"HR9118658081801951861",
	"HR7093391174762888131",
	"HR6824554207539191367",
	"HR7069604594001692768",
	"HR4270163171014341308",
	"CY48590776872388131193566182",
	"CY57511427289148815512463528",
	"CY10469623011747193079305488",
	"CY86826022479357551507194222",
	"CY65139035183710553510799793",
	"CZ3740083988228867633610",
	"CZ4923390395798905071131",
	"CZ3697747307026104738861",
	"CZ5061223246730267064210",
	"CZ3638452806288471256750",
	"DK8387188644726815",
	"DK3068706775436067",
	"DK0697654450063121",
	"DK1099979861456738",
	"DK6988299842527195",
	"EE416702219844182076",
	"EE816382882633746409",
	"EE035815981173988529",
	"EE150595733987082728",
	"EE605409451030627522",
	"FO1593707486505366",
	"FO5006907768039839",
	"FO9378537341306148",
	"FO4068759083981752",
	"FO0905894981715676",
	"FI8709549333658747",
	"FI1518471099159022",
	"FI0589161476500751",
	"FI8433982173935580",
	"FI7954405150189238",
	"FR4197944644738285027717680",
	"FR9476231310567227640169067",
	"FR6888474339535547405026268",
	"FR3007344050937354660134854",
	"FR8547764510959591053030050",
	"PF8169352568136984283973639",
	"PF1021725003919279759512045",
	"PF9067348885442846702112667",
	"PF2110055440192380776287331",
	"PF4462138104308037716665461",
	"TF7369356610212036082878842",
	"TF1699071511365858327828309",
	"TF9287657455706592772258930",
	"TF6983084059527026532259346",
	"TF6320136548014311655407753",
	"DE06495352657836424132",
	"DE09121688720378475751",
	"DE88516399675378845887",
	"DE42187384985716759572",
	"DE04399340668928275395",
	"GI84YQVE742322843673354",
	"GI50TRZE832226672231136",
	"GI96DQBV940980418323607",
	"GI12MTEJ300936244995281",
	"GI50NKEA869461619367593",
	"GE27JX0363248286073573",
	"GE95PE2036699405919650",
	"GE86WI1894058889642409",
	"GE50ZK0956993434292828",
	"GE60VX8276008964044900",
	"GR8206922880502260960449182",
	"GR0708312360607104632724143",
	"GR3019549951345337224826989",
	"GR1850333105485787816165828",
	"GR3328425960116597801941217",
	"GL3357098231928641",
	"GL7672801402871438",
	"GL8576657033000228",
	"GL7533425696727320",
	"GL3145184332080211",
	"HU53165228954563006441576439",
	"HU61442178338678431742505774",
	"HU64774233934011029174507108",
	"HU79689064758089616754511009",
	"HU60873329200412252359645504",
	"IS098954934397185843549690",
	"IS367580035808668402924142",
	"IS179684724271989278617740",
	"IS846240236716627404368872",
	"IS521951362135843206164749",
	"IE49BENI35972029450251",
	"IE37OGUG54280567980573",
	"IE43DCUZ91231044680662",
	"IE15AAKO11199097933967",
	"IE77PIHS49175290558839",
	"IL454322198734138455151",
	"IL839799606658366056087",
	"IL038569613554041996868",
	"IL813919026399312117293",
	"IL654645042217944600527",
	"IT85M5508898545109326040966",
	"IT52G4674641537627600627273",
	"IT54K9621595703270001697912",
	"IT02F7240326523239438656917",
	"IT75F6444007486118207984348",
	"LV85QMUO0600628590552",
	"LV06FYUQ8115346663782",
	"LV05OXNQ0057259369767",
	"LV22XGHZ6356462010762",
	"LV27LLIK8896580861638",
	"LB82586807590631203627574587",
	"LB33405622563828555835796785",
	"LB04715710805951055803616185",
	"LB61420797023022242826619522",
	"LB67864629408749872547678117",
	"LI4091221689235313176",
	"LI7615336074136062084",
	"LI3727301137968672218",
	"LI3551318446915634574",
	"LI4705272204109186337",
	"LT369967216439021801",
	"LT937444321684957069",
	"LT424971109068400772",
	"LT566572547785167976",
	"LT344806290778854389",
	"LU292357816107922497",
	"LU184883493877746720",
	"LU850789684146586224",
	"LU365548629753608759",
	"LU954093702688849179",
	"MK72125600332161582",
	"MK11033425562019483",
	"MK28337919411434742",
	"MK22345789402386151",
	"MK82644233974800672",
	"MT97ATVB58306859106316239974172",
	"MT74VCFO64435204415027820548935",
	"MT35ITGC82712389863518284695353",
	"MT68DQVR03392795978045273628339",
	"MT29SUTJ80635803074721583494800",
	"MU51SJFJ6257989899845328236GLS",
	"MU61KWPF5078030841109086598WUO",
	"MU53JZOY7025842098740945151ZDV",
	"MU33GHPP0512367410476524003SGD",
	"MU47IDNI5979337138037202943JSF",
	"YT2364450161155207655772895",
	"YT5387838908762423789181088",
	"YT5732176546694896615831590",
	"YT9824037454721994164038623",
	"YT2841514046462334743686132",
	"MC7310021462396304214555821",
	"MC7426943447019580313912629",
	"MC5828214954205338889828744",
	"MC5492313829455176982975920",
	"MC9281452662355894512310924",
	"ME60043032533135219382",
	"ME13188638660227646081",
	"ME82608318996043837340",
	"ME15121909794401990976",
	"ME76417412116089156198",
	"NL23CGPQ0251595242",
	"NL21FPBW0870850199",
	"NL15TVWK0331902885",
	"NL13RTEF0518590011",
	"NL40SGFW1252215983",
	"NC9532788614647022310269396",
	"NC2053292379717332255189037",
	"NC9105701404726570049169877",
	"NC1021801619496974025318651",
	"NC1729434258559239166499130",
	"NO5384279272034",
	"NO9739077211102",
	"NO7009234138626",
	"NO4448361377130",
	"NO7962522169141",
	"PL58515427093787930748060666",
	"PL22980511341176988398762666",
	"PL08239642036391620641611736",
	"PL67449258602191338152126294",
	"PL82771306277364889467742211",
	"PT12014625392693045083592",
	"PT80898047569635824751698",
	"PT12065330847682220414039",
	"PT98681708278396096913836",
	"PT49242951581988705914025",
	"RO14JLFB9551925334163469",
	"RO81QBBE5290470985636122",
	"RO11VYHO3215271561449480",
	"RO57EAOP2023783320803714",
	"RO21HNFU2813681045796465",
	"PM8059411251360674293481450",
	"PM5203212193960732379060042",
	"PM7466602890486264340672969",
	"PM6367055534424386034425612",
	"PM2060260873302215070303208",
	"SM70N8724751165335491824812",
	"SM78X1135489836211118891839",
	"SM97M4888143036388138800185",
	"SM72C9584723533916792029340",
	"SM90M9981908196491432695525",
	"SA5591720552379162070001",
	"SA0545544944406431392597",
	"SA1667630781004847967711",
	"SA2589813740329129166910",
	"SA7081962486570441251637",
	"RS85033307149788542871",
	"RS52665698845368481211",
	"RS82691654340096587307",
	"RS88844660406878687897",
	"RS55472917853273859291",
	"SK4167111421162529673536",
	"SK4589732621505695319336",
	"SK4492457066924445710519",
	"SK9190300791649333346556",
	"SK6835978956449243145407",
	"SI26085198624502816",
	"SI85363467889027196",
	"SI93016808632808860",
	"SI45000543512611896",
	"SI14647150971707561",
	"ES2364265841767173822054",
	"ES5577644480024527929849",
	"ES7502766977729557202723",
	"ES3282395478259622275430",
	"ES9034258324029250165663",
	"SE3159169406714737443256",
	"SE2636432381651868407029",
	"SE3280552515152942260664",
	"SE7905464316022155413548",
	"SE8953084170161031273426",
	"CH1987364322975299818",
	"CH4269286867620396437",
	"CH2296292579429731980",
	"CH6518929919723772608",
	"CH9093021641139942126",
	"TN9670288139885457943736",
	"TN8738524364626879391407",
	"TN7275949269046889239714",
	"TN4006837077003057397517",
	"TN8683931110271287238460",
	"TR493318798613751080384953",
	"TR314256533622834177853745",
	"TR080572402207758013538147",
	"TR489116538521358266645727",
	"TR795585070398853758044433",
	"GB39MUJS50172570996370",
	"GB14SIPV86193224493527",
	"GB55ZAFY89851748597528",
	"GB22KVUM18028477988401",
	"GB26JAYK66540091518150",
	"WF5664222423044623595985593",
	"WF6125565335534356679570561",
	"WF4041383920092945092359281",
	"WF0721812715683400832634716",
	"WF6876262234744814330049391"}
