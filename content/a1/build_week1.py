"""Authoring source for A1 Week 1 lessons.

Builds the content bundle imported at /admin/content/import and the audio
manifest used by generate_lesson_audio.sh. Helpers: lesson_kit.py.

Usage (from content/a1):
    python3 build_week1.py [--vocab ../a1-vocabulary/out/a1_rows.json]
Outputs: week-1.json, audio_manifest_week-1.tsv
"""
from lesson_kit import *  # noqa: F401,F403 — exercise / activity helpers
from lesson_kit import finish, lesson, start_week

start_week(
    1,
    # Speaker genders per dialogue (A, B, C…), for the TTS voices.
    dialogue_voices={"l1": "MF", "l2": "MF", "l3": "MFM", "l4": "FM", "l5": "FM",
                     "l6": "MF", "l7": "MF", "l8": "FM", "l9": "FM"},
    reading_genres={"l1": "NOTICE", "l2": "PROFILE", "l3": "EMAIL", "l4": "PROFILE", "l5": "MESSAGE",
                    "l6": "MESSAGE", "l7": "PROFILE", "l8": "NOTICE", "l9": "MESSAGE"},
)

# ---------------------------------------------------------------------------
# Week 1 lessons

G1 = "grammar:g1-to-be"
G2 = "grammar:g2-nouns-articles"
lessons = []

# L1 Greetings ---------------------------------------------------------------
k = lesson("l1")
lessons.append({
    "slug": "a1-w1-l1-greetings", "title": "Greetings",
    "description": "Chào hỏi theo thời điểm trong ngày, chào tạm biệt, cảm ơn và xin lỗi.",
    "activities": [
        pronunciation("Âm /iː/ và /ɪ/", "/iː/ là âm i dài (kéo dài, miệng cười), /ɪ/ là âm i ngắn và thả lỏng.", ["iː", "ɪ"], [
            minimal_pair(k + "1", [("ship", "/ʃɪp/"), ("sheep", "/ʃiːp/")], 1, "'sheep' có /iː/ dài, 'ship' có /ɪ/ ngắn.", "sound:iː"),
            minimal_pair(k + "2", [("live", "/lɪv/"), ("leave", "/liːv/")], 0, "'live' (sống) có /ɪ/ ngắn, 'leave' (rời đi) có /iː/ dài.", "sound:ɪ"),
            ipa_to_word("/miːt/", ["meet", "mitt", "met"], 0, "/iː/ thường được viết là 'ee' hoặc 'ea'.", "sound:iː"),
            stress(k, "hello", ["hel", "lo"], 1, "hel-LO: nhấn âm tiết thứ hai."),
        ]),
        vocabulary("Từ vựng chào hỏi", ["a1-hello", "a1-hi", "a1-good-morning", "a1-good-afternoon", "a1-good-night",
                                        "a1-bye", "a1-welcome", "a1-please", "a1-thank", "a1-sorry"], [
            mc("\"Good afternoon\" dùng khi nào?", ["Sau 12 giờ trưa", "Trước 12 giờ trưa", "Trước khi đi ngủ"], 0, "Good afternoon dùng từ 12 giờ trưa đến chiều tối.", "vocab:greetings"),
            mc("Câu nào dùng trước khi đi ngủ?", ["Good morning", "Good night", "Good afternoon"], 1, "Good night = chúc ngủ ngon, chỉ dùng khi chia tay buổi tối hoặc đi ngủ.", "vocab:greetings"),
            mc("Bạn đến muộn, bạn nói gì?", ["Thank you.", "Sorry, I'm late.", "Welcome!"], 1, "Sorry = xin lỗi.", "vocab:greetings"),
        ]),
        grammar("Động từ to be: am / is / are", "I đi với am; he, she, it đi với is; you, we, they đi với are.", "g1-to-be", [
            fill("I ___ (be) Huy.", ["am"], "Chủ ngữ 'I' luôn đi với 'am'.", G1),
            fill("She ___ (be) my teacher.", ["is"], "He / She / It đi với 'is'.", G1),
            fill("We ___ (be) colleagues.", ["are"], "We / You / They đi với 'are'.", G1),
            reorder("How are you?", "Câu hỏi thăm sức khỏe: How are you?", G1),
            error_correction("He are fine.", 1, ["is"], "'He' đi với 'is'.", G1),
        ]),
        listening(k, "Chào buổi sáng ở văn phòng", [
            ("A", "Good morning, Lan."), ("B", "Good morning, Minh. How are you?"),
            ("A", "I'm fine, thank you. And you?"), ("B", "I'm very well, thanks."),
            ("A", "See you at the meeting. Bye!"), ("B", "Bye!"),
        ], [
            true_false("It is the morning.", True, "Họ nói 'Good morning'.", "listening"),
            true_false("Lan is not well.", False, "Lan nói 'I'm very well' (tôi rất khỏe).", "listening"),
            dictation(k + "1", "I'm fine, thank you.", "I'm fine = tôi khỏe."),
        ]),
        reading("Thông báo đầu tuần", "Good morning, everyone! Welcome back to work. The team meeting is at 9 o'clock in Room A. Please be on time. Have a nice day! Anna, Office Manager", [
            mc("Cuộc họp lúc mấy giờ?", ["9 o'clock", "10 o'clock", "8 o'clock"], 0, "\"The team meeting is at 9 o'clock.\"", "reading"),
            true_false("The meeting is in Room B.", False, "Cuộc họp ở Room A.", "reading"),
            true_false("Anna is the office manager.", True, "\"Anna, Office Manager\"", "reading"),
            mc("Thông báo này viết cho ai?", ["everyone", "Anna", "the manager"], 0, "\"Good morning, everyone!\"", "reading"),
        ]),
        speaking(["Good morning! How are you?", "I'm fine, thank you. And you?", "Sorry, I'm late."]),
    ],
})

# L2 Introducing yourself ----------------------------------------------------
k = lesson("l2")
lessons.append({
    "slug": "a1-w1-l2-introducing-yourself", "title": "Introducing yourself",
    "description": "Nói tên, công việc, nơi sống và quê quán.",
    "activities": [
        pronunciation("Âm /e/ và /æ/", "/e/ miệng mở vừa (như 'e' tiếng Việt); /æ/ miệng mở rộng hơn, giữa 'e' và 'a'.", ["e", "æ"], [
            minimal_pair(k + "1", [("men", "/men/"), ("man", "/mæn/")], 1, "'man' có /æ/ (mở miệng rộng), 'men' có /e/.", "sound:æ"),
            minimal_pair(k + "2", [("bed", "/bed/"), ("bad", "/bæd/")], 0, "'bed' có /e/, 'bad' có /æ/.", "sound:e"),
            word_to_ipa(k + "1", "name", ["/neɪm/", "/næm/", "/nem/"], 0, "'name' có nguyên âm đôi /eɪ/.", "sound:eɪ"),
            ipa_to_word("/hæpi/", ["happy", "hippie", "help"], 0, "/æ/ trong 'happy'.", "sound:æ"),
        ]),
        vocabulary("Thông tin cá nhân", ["a1-name", "a1-last-name", "a1-live", "a1-from", "a1-hometown",
                                          "a1-nationality", "a1-job", "a1-work", "a1-meet", "a1-introduce"], [
            mc("\"Hometown\" nghĩa là gì?", ["quê hương", "thành phố lớn", "nhà riêng"], 0, "Hometown = nơi bạn sinh ra / quê.", "vocab:personal"),
            mc("Chọn câu đúng để nói nơi bạn sống", ["I live in Hanoi.", "I from Hanoi.", "I am live in Hanoi."], 0, "live in + nơi chốn.", "vocab:personal"),
            mc("\"Nice to meet you\" dùng khi nào?", ["Lần đầu gặp ai đó", "Khi tạm biệt buổi tối", "Khi xin lỗi"], 0, "Nice to meet you = rất vui được gặp bạn.", "vocab:personal"),
        ]),
        grammar("I am / I'm — câu khẳng định", "Trong văn nói thường dùng dạng rút gọn: I'm, you're, he's, she's, we're, they're.", "g1-to-be", [
            fill("My name ___ Huy.", ["is"], "'My name' là ngôi thứ ba số ít → is.", G1),
            fill("I ___ (be) from Vietnam.", ["am", "'m"], "I am from + quốc gia.", G1),
            reorder("I am a developer.", "Nghề nghiệp: I am + a/an + nghề.", G1),
            translate("Tôi sống ở Hà Nội.", ["I live in Hanoi.", "I live in Ha Noi."], "live in + thành phố.", G1),
            translate("Tôi là kỹ sư.", ["I am an engineer.", "I'm an engineer."], "engineer bắt đầu bằng nguyên âm → dùng 'an'.", G1),
        ]),
        listening(k, "Ngày đầu đi làm", [
            ("A", "Hi, I'm Minh. I'm the new designer."), ("B", "Nice to meet you, Minh. I'm Lan."),
            ("A", "Nice to meet you too. Where are you from?"), ("B", "I'm from Da Nang, but I live in Hanoi now."),
        ], [
            mc("Minh làm nghề gì?", ["designer", "developer", "teacher"], 0, "\"I'm the new designer.\"", "listening"),
            true_false("Lan lives in Da Nang now.", False, "Lan quê Đà Nẵng nhưng giờ sống ở Hà Nội.", "listening"),
            dictation(k + "1", "Nice to meet you, Minh.", "Câu chào khi lần đầu gặp."),
        ]),
        reading("Hồ sơ nhân viên mới", "Hello! My name is Mai. I'm 26 years old. I'm from Hue, but I live in Ho Chi Minh City. I'm an accountant at a tech company. Nice to meet you!", [
            true_false("Mai is from Ho Chi Minh City.", False, "Mai quê ở Huế.", "reading"),
            mc("Mai bao nhiêu tuổi?", ["26", "36", "16"], 0, "\"I'm 26 years old.\"", "reading"),
            mc("Mai làm nghề gì?", ["accountant", "teacher", "designer"], 0, "\"I'm an accountant.\"", "reading"),
            true_false("Mai lives in Ho Chi Minh City.", True, "\"...but I live in Ho Chi Minh City.\"", "reading"),
        ]),
        speaking(["Hi, I'm Huy. I'm a developer.", "I'm from Vietnam. I live in Hanoi.", "Nice to meet you!"]),
    ],
})

# L3 Introducing others ------------------------------------------------------
k = lesson("l3")
lessons.append({
    "slug": "a1-w1-l3-introducing-others", "title": "Introducing others",
    "description": "Giới thiệu đồng nghiệp, bạn bè: This is…, He is…, She is…",
    "activities": [
        pronunciation("Âm /ʌ/ và /ɑː/", "/ʌ/ ngắn, miệng mở vừa (như 'ơ' ngắn); /ɑː/ dài, miệng mở rộng (như 'a' kéo dài).", ["ʌ", "ɑː"], [
            minimal_pair(k + "1", [("cut", "/kʌt/"), ("cart", "/kɑːrt/")], 0, "'cut' có /ʌ/ ngắn.", "sound:ʌ"),
            minimal_pair(k + "2", [("hut", "/hʌt/"), ("hot", "/hɑːt/")], 1, "'hot' có /ɑː/ trong giọng Mỹ.", "sound:ɑː"),
            ipa_to_word("/brʌðər/", ["brother", "bother", "breather"], 0, "'brother' có /ʌ/.", "sound:ʌ"),
            stress(k, "colleague", ["col", "league"], 0, "COL-league: nhấn âm tiết đầu."),
        ]),
        vocabulary("Người xung quanh", ["a1-friend", "a1-colleague", "a1-boss", "a1-manager", "a1-team",
                                         "a1-man", "a1-woman", "a1-he", "a1-she", "a1-they"], [
            mc("\"Colleague\" nghĩa là gì?", ["đồng nghiệp", "sếp", "khách hàng"], 0, "Colleague = người làm cùng chỗ với bạn.", "vocab:people"),
            mc("Chọn từ chỉ người quản lý", ["manager", "customer", "guest"], 0, "Manager = quản lý.", "vocab:people"),
            mc("Giới thiệu một người phụ nữ, bạn dùng đại từ nào?", ["she", "he", "it"], 0, "She dùng cho nữ, he dùng cho nam.", "vocab:people"),
        ]),
        grammar("This is… / He is… / She is…", "Giới thiệu người khác: This is + tên. Sau đó dùng He is / She is để nói thêm.", "g1-to-be", [
            fill("This ___ my friend, Nam.", ["is"], "This is + tên người.", G1),
            fill("They ___ (be) my colleagues.", ["are"], "They đi với are.", G1),
            reorder("She is our new manager.", "Cấu trúc: She is + a/our + chức danh.", G1),
            error_correction("Lan and Minh is designers.", 3, ["are"], "Hai người (Lan and Minh) → are.", G1),
            transform("He is a teacher.", "THIRD_PERSON", ["She is a teacher.", "She's a teacher."], "Đổi He thành She, động từ vẫn là is.", G1),
        ]),
        listening(k, "Giới thiệu đồng nghiệp", [
            ("A", "Lan, this is Peter. He is our new developer."), ("B", "Hi, Peter. Nice to meet you."),
            ("C", "Nice to meet you too, Lan."), ("A", "Peter is from Canada. He is very friendly."),
        ], [
            true_false("Peter is a new developer.", True, "\"He is our new developer.\"", "listening"),
            mc("Peter đến từ đâu?", ["Canada", "Australia", "the USA"], 0, "\"Peter is from Canada.\"", "listening"),
            dictation(k + "1", "This is Peter.", "This is + tên: giới thiệu ai đó."),
        ]),
        reading("Email giới thiệu người mới", "Hi team, Please welcome Peter Brown. He is our new developer. He is from Toronto, Canada, and he is 29 years old. Peter is in the tech team on the third floor. Please say hello to him! Thanks, Lan", [
            true_false("Peter is a designer.", False, "\"He is our new developer.\"", "reading"),
            mc("Peter đến từ thành phố nào?", ["Toronto", "Sydney", "London"], 0, "\"He is from Toronto, Canada.\"", "reading"),
            mc("Peter bao nhiêu tuổi?", ["29", "19", "39"], 0, "\"he is 29 years old\"", "reading"),
            true_false("The tech team is on the third floor.", True, "\"in the tech team on the third floor\"", "reading"),
        ]),
        speaking(["This is my colleague, Lan.", "She is a designer.", "They are my friends."]),
    ],
})

# L4 Personal information -----------------------------------------------------
k = lesson("l4")
lessons.append({
    "slug": "a1-w1-l4-personal-information", "title": "Asking for personal information",
    "description": "Hỏi và trả lời về tên, tuổi, số điện thoại, email: What's your…?",
    "activities": [
        pronunciation("Âm /uː/ và /ʊ/", "/uː/ dài, môi tròn và chu ra; /ʊ/ ngắn, môi thả lỏng hơn.", ["uː", "ʊ"], [
            minimal_pair(k + "1", [("full", "/fʊl/"), ("fool", "/fuːl/")], 1, "'fool' có /uː/ dài.", "sound:uː"),
            minimal_pair(k + "2", [("look", "/lʊk/"), ("Luke", "/luːk/")], 0, "'look' có /ʊ/ ngắn.", "sound:ʊ"),
            ipa_to_word("/nuːn/", ["noon", "nun", "none"], 0, "/uː/ thường viết là 'oo'.", "sound:uː"),
            stress(k, "address", ["ad", "dress"], 0, "Giọng Mỹ: AD-dress (danh từ)."),
        ]),
        vocabulary("Thông tin liên lạc", ["a1-name", "a1-age", "a1-old", "a1-address", "a1-phone",
                                           "a1-email", "a1-number", "a1-married", "a1-single", "a1-what"], [
            mc("\"Address\" nghĩa là gì?", ["địa chỉ", "số điện thoại", "tuổi"], 0, "Address = địa chỉ.", "vocab:personal"),
            mc("Hỏi tuổi đúng là câu nào?", ["How old are you?", "How age are you?", "What old are you?"], 0, "How old are you? = Bạn bao nhiêu tuổi?", "vocab:personal"),
            mc("\"Single\" nghĩa là gì?", ["độc thân", "đã kết hôn", "một mình"], 0, "Single = độc thân (chưa kết hôn).", "vocab:personal"),
        ]),
        grammar("Câu hỏi với to be", "Đảo to be lên trước chủ ngữ: Are you…? Is she…? Hoặc dùng từ để hỏi: What is your…?", "g1-to-be", [
            fill("___ you married?", ["are"], "Câu hỏi Yes/No: Are + you + …?", G1),
            fill("What ___ your phone number?", ["is"], "What is your + danh từ số ít?", G1),
            transform("She is single.", "QUESTION", ["Is she single?"], "Đưa 'is' lên trước chủ ngữ.", G1),
            transform("You are 30 years old.", "NEGATIVE", ["You are not 30 years old.", "You aren't 30 years old.", "You're not 30 years old."], "Phủ định: are + not.", G1),
            reorder("What is your email address?", "What is your + thông tin?", G1),
        ]),
        listening(k, "Điền thông tin ở lễ tân", [
            ("A", "What's your name, please?"), ("B", "My name is Huy. H-U-Y."),
            ("A", "And what's your phone number?"), ("B", "It's 0912 345 678."),
            ("A", "How old are you?"), ("B", "I'm thirty-two."),
        ], [
            numbers(k + "1", "It's 0912 345 678.", "Số điện thoại của Huy là gì?", "PHONE", ["0912 345 678"], "Đọc từng số một."),
            numbers(k + "2", "I'm thirty-two.", "Huy bao nhiêu tuổi?", "NUMBER", ["32"], "thirty-two = 32."),
            true_false("Huy spells his name H-U-Y.", True, "Đánh vần: H-U-Y.", "listening"),
        ]),
        reading("Phiếu thông tin nhân viên", "Name: Nguyen Van Huy\nAge: 32\nPhone: 0912 345 678\nEmail: huy.nv@example.com\nAddress: 15 Tran Hung Dao Street, Hanoi\nJob: software developer", [
            mc("Huy sống ở thành phố nào?", ["Hanoi", "Hue", "Da Nang"], 0, "Address: ... Hanoi.", "reading"),
            true_false("Huy is 23 years old.", False, "Age: 32 (thirty-two), không phải 23.", "reading"),
            mc("Huy làm nghề gì?", ["software developer", "doctor", "teacher"], 0, "Job: software developer.", "reading"),
            mc("Số nhà của Huy là bao nhiêu?", ["15", "50", "5"], 0, "15 Tran Hung Dao Street.", "reading"),
        ]),
        speaking(["What's your name?", "What's your phone number?", "How old are you?"]),
    ],
})

# L5 Family -------------------------------------------------------------------
k = lesson("l5")
lessons.append({
    "slug": "a1-w1-l5-family", "title": "Talking about your family",
    "description": "Nói về các thành viên trong gia đình: I have…, my father, my sister…",
    "activities": [
        pronunciation("Âm /ɜːr/ và /ə/", "/ɜːr/ dài, lưỡi cong (như trong 'work', 'girl'); /ə/ là âm 'ơ' rất ngắn, nhẹ, không nhấn.", ["ɜːr", "ə"], [
            minimal_pair(k + "1", [("bird", "/bɜːrd/"), ("bed", "/bed/")], 0, "'bird' có /ɜːr/.", "sound:ɜːr"),
            ipa_to_word("/ɡɜːrl/", ["girl", "goal", "gull"], 0, "'girl' có /ɜːr/.", "sound:ɜːr"),
            ipa_to_word("/ˈmʌðər/", ["mother", "matter", "mutter"], 0, "Âm cuối '-er' đọc /ər/ nhẹ.", "sound:ə"),
            stress(k, "parents", ["par", "ents"], 0, "PAR-ents."),
        ]),
        vocabulary("Gia đình", ["a1-family", "a1-father", "a1-mother", "a1-parent", "a1-brother", "a1-sister",
                                "a1-son", "a1-daughter", "a1-husband", "a1-wife"], [
            mc("\"Daughter\" nghĩa là gì?", ["con gái", "con trai", "chị gái"], 0, "Daughter = con gái; son = con trai.", "vocab:family"),
            mc("Bố và mẹ gọi chung là gì?", ["parents", "children", "cousins"], 0, "Parents = bố mẹ.", "vocab:family"),
            mc("\"Husband\" là gì?", ["chồng", "vợ", "anh trai"], 0, "Husband = chồng; wife = vợ.", "vocab:family"),
        ]),
        grammar("Danh từ số nhiều và a/an", "Số nhiều thường thêm -s (sisters) hoặc -es (boxes). Dùng a trước phụ âm, an trước nguyên âm. Một số từ bất quy tắc: child → children, man → men, woman → women.", "g2-nouns-articles", [
            fill("I have two ___ (sister).", ["sisters"], "two + danh từ số nhiều: sister → sisters.", G2),
            fill("They have three ___ (child).", ["children"], "child → children (bất quy tắc).", G2),
            mc("Chọn mạo từ đúng: ___ engineer", ["an", "a", "the"], 0, "engineer bắt đầu bằng nguyên âm /e/ → an.", G2),
            error_correction("My brother has two son.", 4, ["sons"], "two → danh từ số nhiều: sons.", G2),
            translate("Tôi có một anh trai và hai em gái.", ["I have a brother and two sisters.", "I have one brother and two sisters.", "I have an older brother and two younger sisters."], "have + a/one + danh từ; two + danh từ số nhiều.", G2),
        ]),
        listening(k, "Ảnh gia đình", [
            ("A", "Is this your family?"), ("B", "Yes. This is my father and this is my mother."),
            ("A", "Who is this?"), ("B", "That's my sister, Hoa. She is a nurse. And these are her two children."),
        ], [
            mc("Hoa làm nghề gì?", ["nurse", "doctor", "teacher"], 0, "\"She is a nurse.\"", "listening"),
            numbers(k + "1", "These are her two children.", "Hoa có mấy con?", "NUMBER", ["2"], "two children = 2 đứa con."),
            dictation(k + "1", "This is my father.", ""),
        ]),
        reading("Tin nhắn của Nam", "Hi Lan! Thanks for the invitation. I can come to the party with my wife, Linh. Our son is five and our daughter is two. Is that OK? Nam", [
            true_false("Nam has two children.", True, "Một con trai và một con gái.", "reading"),
            mc("Con trai Nam bao nhiêu tuổi?", ["5", "2", "7"], 0, "\"Our son is five.\"", "reading"),
            mc("Vợ Nam tên gì?", ["Linh", "Lan", "Hoa"], 0, "\"my wife, Linh\"", "reading"),
            true_false("Nam writes to Lan.", True, "Tin nhắn bắt đầu bằng \"Hi Lan!\".", "reading"),
        ]),
        speaking(["I have one brother and two sisters.", "My mother is a teacher.", "This is my family."]),
    ],
})

# L6 Hobbies ------------------------------------------------------------------
k = lesson("l6")
lessons.append({
    "slug": "a1-w1-l6-hobbies", "title": "Talking about hobbies",
    "description": "Nói về sở thích: I like…, My hobby is…, đọc đuôi -s.",
    "activities": [
        pronunciation("Đọc đuôi -s: /s/, /z/, /ɪz/", "Sau âm vô thanh (p, t, k, f) đọc /s/; sau âm hữu thanh và nguyên âm đọc /z/; sau s, z, sh, ch, x đọc /ɪz/.", [], [
            ending(k + "1", "books", "/s/", "Sau /k/ (vô thanh) → /s/."),
            ending(k + "2", "songs", "/z/", "Sau /ŋ/ (hữu thanh) → /z/."),
            ending(k + "3", "watches", "/ɪz/", "Sau 'ch' → /ɪz/."),
            ending(k + "4", "plays", "/z/", "Sau nguyên âm → /z/."),
        ]),
        vocabulary("Sở thích", ["a1-hobby", "a1-free-time", "a1-enjoy", "a1-play", "a1-dance", "a1-sing",
                                 "a1-draw", "a1-fishing", "a1-music", "a1-movie"], [
            mc("\"Free time\" nghĩa là gì?", ["thời gian rảnh", "miễn phí", "thời gian làm việc"], 0, "Free time = thời gian rảnh.", "vocab:hobbies"),
            mc("Chọn hoạt động đúng với 'fishing'", ["câu cá", "nấu ăn", "bơi"], 0, "Fishing = câu cá.", "vocab:hobbies"),
            mc("\"I enjoy cooking\" nghĩa là?", ["Tôi thích nấu ăn", "Tôi biết nấu ăn", "Tôi đang nấu ăn"], 0, "enjoy + V-ing = thích làm gì.", "vocab:hobbies"),
        ]),
        grammar("Ôn to be + sở hữu my / your", "My hobby is… (sở thích của tôi là…). Your = của bạn, his = của anh ấy, her = của cô ấy.", "g1-to-be", [
            fill("My hobby ___ reading.", ["is"], "My hobby (số ít) → is.", G1),
            fill("What ___ your hobbies?", ["are"], "your hobbies (số nhiều) → are.", G1),
            mc("Chọn từ đúng: This is Lan. ___ hobby is painting.", ["Her", "His", "Its"], 0, "Lan là nữ → her.", G2),
            reorder("My favorite hobby is music.", "favorite = yêu thích nhất.", G1),
            error_correction("His hobbies is football and music.", 2, ["are"], "Hobbies (số nhiều) → are.", G1),
        ]),
        listening(k, "Cuối tuần bạn làm gì?", [
            ("A", "What do you do in your free time, Lan?"), ("B", "I like music. I sing and play the guitar."),
            ("A", "Cool! My hobby is fishing. I go fishing with my father on Sunday."),
        ], [
            mc("Lan chơi nhạc cụ gì?", ["guitar", "piano", "drums"], 0, "\"I sing and play the guitar.\"", "listening"),
            true_false("The man goes fishing on Saturday.", False, "Anh ấy đi câu cá vào Chủ nhật (Sunday).", "listening"),
            gap_fill(k + "1", "I go fishing with my father on Sunday.", "I go {{1}} with my father on {{2}}.",
                     {"1": ["fishing"], "2": ["Sunday"]}),
        ]),
        reading("Rủ đi lễ hội âm nhạc", "Hi Minh! Are you free on Saturday? There is a music festival in the park. My sister and I love rock music. We have three tickets. Do you want to come with us? Lan", [
            mc("Lễ hội âm nhạc ở đâu?", ["in the park", "at school", "at the office"], 0, "\"There is a music festival in the park.\"", "reading"),
            true_false("Lan likes rock music.", True, "\"My sister and I love rock music.\"", "reading"),
            mc("Họ có mấy vé?", ["3", "2", "4"], 0, "\"We have three tickets.\"", "reading"),
            true_false("The festival is on Sunday.", False, "Lễ hội vào thứ Bảy (Saturday).", "reading"),
        ]),
        speaking(["My hobby is reading.", "I like music and movies.", "What do you do in your free time?"]),
    ],
})

# L7 Describing your job --------------------------------------------------------
k = lesson("l7")
lessons.append({
    "slug": "a1-w1-l7-describing-your-job", "title": "Describing your job",
    "description": "Nói về nghề nghiệp và nơi làm việc: I'm a/an…, I work for/in…",
    "activities": [
        pronunciation("Trọng âm từ chỉ nghề nghiệp", "Từ hai âm tiết chỉ người thường nhấn âm tiết đầu (TEACH-er), nhưng có ngoại lệ (en-gi-NEER).", [], [
            stress(k + "1", "teacher", ["teach", "er"], 0, "TEACH-er."),
            stress(k + "2", "doctor", ["doc", "tor"], 0, "DOC-tor."),
            stress(k + "3", "designer", ["de", "sign", "er"], 1, "de-SIGN-er."),
            stress(k + "4", "engineer", ["en", "gi", "neer"], 2, "en-gi-NEER: nhấn âm cuối."),
        ]),
        vocabulary("Nghề nghiệp", ["a1-teacher", "a1-doctor", "a1-nurse", "a1-engineer", "a1-driver", "a1-developer",
                                    "a1-designer", "a1-accountant", "a1-company", "a1-office"], [
            mc("\"Accountant\" là nghề gì?", ["kế toán", "kỹ sư", "lập trình viên"], 0, "Accountant = kế toán.", "vocab:jobs"),
            mc("Người lái xe taxi gọi là gì?", ["driver", "doctor", "designer"], 0, "Driver = tài xế.", "vocab:jobs"),
            mc("Chọn câu đúng", ["I'm an accountant.", "I'm a accountant.", "I'm accountant."], 0, "accountant bắt đầu bằng nguyên âm → an.", G2),
        ]),
        grammar("a / an trước nghề nghiệp", "Khi nói nghề, luôn có a/an: I'm a teacher. I'm an engineer. Số nhiều không dùng a/an: They are nurses.", "g2-nouns-articles", [
            fill("She is ___ nurse.", ["a"], "nurse bắt đầu bằng phụ âm → a.", G2),
            fill("He is ___ engineer.", ["an"], "engineer bắt đầu bằng nguyên âm → an.", G2),
            error_correction("They are a doctors.", 2, ["doctors"], "Số nhiều không dùng 'a'. Sửa thành 'They are doctors.'", G2),
            translate("Tôi là lập trình viên ở một công ty công nghệ.", ["I am a developer at a tech company.", "I'm a developer at a tech company.", "I am a developer in a tech company.", "I'm a developer in a technology company.", "I am a developer at a technology company."], "a developer at a + công ty.", G2),
            reorder("I work for a big company.", "work for + công ty.", G2),
        ]),
        listening(k, "Bạn làm nghề gì?", [
            ("A", "What do you do, Hoa?"), ("B", "I'm an accountant. I work for a bank."),
            ("A", "Is it a big bank?"), ("B", "Yes, it is. There are two hundred people in my office."),
        ], [
            mc("Hoa làm việc ở đâu?", ["a bank", "a hospital", "a school"], 0, "\"I work for a bank.\"", "listening"),
            numbers(k + "1", "There are two hundred people in my office.", "Văn phòng của Hoa có bao nhiêu người?", "NUMBER", ["200"], "two hundred = 200."),
            dictation(k + "1", "I'm an accountant.", "an + accountant."),
        ]),
        reading("Một ngày của Sarah", "Sarah is a nurse. Her hospital is in the city centre. Her work day is from 7 a.m. to 3 p.m. Her job is hard, but she is happy. Her husband, Tom, is a chef. His restaurant is next to the hospital.", [
            mc("Sarah làm nghề gì?", ["nurse", "chef", "doctor"], 0, "\"Sarah is a nurse.\"", "reading"),
            true_false("Tom is a doctor.", False, "Tom là đầu bếp (chef).", "reading"),
            mc("Sarah làm việc đến mấy giờ?", ["3 p.m.", "7 a.m.", "5 p.m."], 0, "\"from 7 a.m. to 3 p.m.\"", "reading"),
            true_false("Tom's restaurant is next to the hospital.", True, "\"His restaurant is next to the hospital.\"", "reading"),
        ]),
        speaking(["What do you do?", "I'm a developer. I work for a tech company.", "My office is in District 1."]),
    ],
})

# L8 Finding places in the building ---------------------------------------------
k = lesson("l8")
lessons.append({
    "slug": "a1-w1-l8-places-in-the-building", "title": "Finding places in the building",
    "description": "Hỏi phòng ở đâu, tầng mấy: Where is…? It's on the second floor.",
    "activities": [
        pronunciation("Số thứ tự và trọng âm số", "thirteen (thir-TEEN) nhấn sau, thirty (THIR-ty) nhấn trước — dễ nhầm khi nghe.", [], [
            minimal_pair(k + "1", [("thirteen", "/θɜːrˈtiːn/"), ("thirty", "/ˈθɜːrti/")], 0, "thir-TEEN nhấn âm sau, kéo dài /iːn/.", "stress:numbers"),
            minimal_pair(k + "2", [("fourteen", "/fɔːrˈtiːn/"), ("forty", "/ˈfɔːrti/")], 1, "FOR-ty nhấn âm đầu.", "stress:numbers"),
            stress(k + "1", "fifteen", ["fif", "teen"], 1, "fif-TEEN."),
            word_to_ipa(k + "1", "floor", ["/flɔːr/", "/fluːr/", "/flaʊr/"], 0, "floor = /flɔːr/.", "sound:ɔː"),
        ]),
        vocabulary("Trong tòa nhà văn phòng", ["a1-building", "a1-floor", "a1-room", "a1-toilet", "a1-stairs",
                                               "a1-door", "a1-next-to", "a1-near", "a1-where", "a1-second"], [
            mc("\"Stairs\" nghĩa là gì?", ["cầu thang", "thang máy", "hành lang"], 0, "Stairs = cầu thang bộ.", "vocab:building"),
            mc("Tầng 2 nói thế nào?", ["the second floor", "the two floor", "the floor two"], 0, "Dùng số thứ tự: the second floor.", "vocab:building"),
            mc("\"Next to\" nghĩa là gì?", ["bên cạnh", "phía trên", "đối diện"], 0, "Next to = bên cạnh.", "vocab:building"),
        ]),
        grammar("There is / There are + giới từ", "There is + danh từ số ít; there are + danh từ số nhiều. Vị trí: on the … floor, next to, near, in.", "g2-nouns-articles", [
            fill("There ___ a toilet on this floor.", ["is"], "a toilet (số ít) → there is.", G2),
            fill("There ___ two meeting rooms next to the kitchen.", ["are"], "two meeting rooms (số nhiều) → there are.", G2),
            mc("Chọn giới từ đúng: My office is ___ the third floor.", ["on", "in", "at"], 0, "on + the … floor.", G2),
            reorder("Where is the meeting room?", "Where is + the + nơi chốn?", G2),
            transform("There is a printer here.", "QUESTION", ["Is there a printer here?"], "Câu hỏi: Is there + a + danh từ?", G2),
        ]),
        listening(k, "Hỏi đường trong tòa nhà", [
            ("A", "Excuse me, where is the meeting room?"), ("B", "It's on the thirteenth floor, next to the kitchen."),
            ("A", "And where are the toilets?"), ("B", "They're near the stairs."),
        ], [
            mc("Phòng họp ở tầng mấy?", ["13", "30", "3"], 0, "thirteenth = thứ 13.", "listening"),
            true_false("The toilets are near the stairs.", True, "\"They're near the stairs.\"", "listening"),
            gap_fill(k + "1", "It's on the thirteenth floor, next to the kitchen.",
                     "It's on the {{1}} floor, next to the {{2}}.", {"1": ["thirteenth", "13th"], "2": ["kitchen"]}),
        ]),
        reading("Bảng chỉ dẫn tòa nhà", "Floor 1: Reception, Cafe. Floor 2: Sales team, Meeting room A. Floor 3: Tech team, Meeting room B. Toilets are on every floor, next to the stairs.", [
            mc("Meeting room B ở tầng mấy?", ["3", "2", "1"], 0, "Floor 3: Tech team, Meeting room B.", "reading"),
            true_false("There is a cafe on the second floor.", False, "Cafe ở tầng 1.", "reading"),
            true_false("There are toilets on every floor.", True, "\"Toilets are on every floor.\"", "reading"),
            mc("Sales team ở tầng mấy?", ["2", "1", "3"], 0, "Floor 2: Sales team.", "reading"),
        ]),
        speaking(["Excuse me, where is the meeting room?", "It's on the second floor.", "The toilets are next to the stairs."]),
    ],
})

# L9 Borrowing things -----------------------------------------------------------
k = lesson("l9")
lessons.append({
    "slug": "a1-w1-l9-borrowing-things", "title": "Borrowing things at work",
    "description": "Nhờ và mượn đồ lịch sự: Can I borrow…? Can you help me? this / that / these / those.",
    "activities": [
        pronunciation("Ôn tập: /iː/–/ɪ/ và đuôi -s", "Ôn lại âm dài/ngắn và cách đọc đuôi -s của tuần này.", ["iː", "ɪ"], [
            minimal_pair(k + "1", [("these", "/ðiːz/"), ("this", "/ðɪs/")], 0, "'these' có /iː/ dài và kết thúc bằng /z/.", "sound:iː"),
            ending(k + "1", "pens", "/z/", "Sau /n/ (hữu thanh) → /z/."),
            ending(k + "2", "laptops", "/s/", "Sau /p/ (vô thanh) → /s/."),
            ending(k + "3", "boxes", "/ɪz/", "Sau 'x' → /ɪz/."),
        ]),
        vocabulary("Đồ dùng văn phòng", ["a1-pen", "a1-pencil", "a1-notebook", "a1-laptop", "a1-paper",
                                          "a1-printer", "a1-desk", "a1-help", "a1-please", "a1-this"], [
            mc("\"Notebook\" là gì?", ["sổ tay", "máy tính xách tay", "giấy"], 0, "Notebook = sổ tay; laptop = máy tính xách tay.", "vocab:office"),
            mc("Câu nhờ lịch sự nhất là?", ["Can you help me, please?", "Help me.", "You help me."], 0, "Can you … please? là cách nhờ lịch sự.", "vocab:office"),
            mc("\"Printer\" nghĩa là gì?", ["máy in", "máy tính", "bàn làm việc"], 0, "Printer = máy in.", "vocab:office"),
        ]),
        grammar("this / that / these / those", "this (cái này, gần, số ít), that (cái kia, xa, số ít), these (những cái này), those (những cái kia).", "g2-nouns-articles", [
            mc("___ pens here are mine.", ["These", "This", "That"], 0, "pens (số nhiều, ở gần) → these.", G2),
            mc("Is ___ your laptop over there?", ["that", "these", "this"], 0, "over there (ở xa, số ít) → that.", G2),
            error_correction("This books are new.", 0, ["These"], "books là số nhiều → These.", G2),
            reorder("Can I use your pen, please?", "Can I + động từ + …, please?", G2),
            translate("Bạn có thể giúp tôi không?", ["Can you help me?", "Could you help me?", "Can you help me, please?"], "Can you + động từ?", G2),
        ]),
        listening(k, "Mượn đồ đồng nghiệp", [
            ("A", "Excuse me, Nam. Can I use your pen, please?"), ("B", "Sure. Here you are."),
            ("A", "Thanks. And are those your notebooks?"), ("B", "Yes, they are. Do you need one?"),
            ("A", "Yes, please. Thank you so much!"),
        ], [
            mc("Người phụ nữ muốn dùng gì trước tiên?", ["a pen", "a laptop", "a printer"], 0, "\"Can I use your pen?\"", "listening"),
            true_false("The notebooks are Nam's.", True, "\"Are those your notebooks?\" — \"Yes, they are.\"", "listening"),
            dictation(k + "1", "Can I use your pen, please?", ""),
        ]),
        reading("Mẩu giấy nhắn của Mai", "Hi Nam, sorry! I have your stapler and two of your pens. They are on my desk, next to the printer. Can I use your laptop charger this afternoon? Thanks! Mai", [
            true_false("Mai has three of Nam's pens.", False, "\"two of your pens\" — hai cây bút.", "reading"),
            mc("Mai muốn mượn gì chiều nay?", ["a laptop charger", "a stapler", "a printer"], 0, "\"Can I use your laptop charger?\"", "reading"),
            true_false("The pens are next to the printer.", True, "\"They are on my desk, next to the printer.\"", "reading"),
            mc("Ai viết mẩu giấy này?", ["Mai", "Nam", "Lan"], 0, "Ký tên: Mai.", "reading"),
        ]),
        speaking(["Can I use your laptop, please?", "Can you help me, please?", "Thank you so much!"]),
    ],
})

# Checkpoint lesson (week review) --------------------------------------------------
k = lesson("rv")
lessons.append({
    "slug": "a1-w1-review", "title": "Week 1 review",
    "description": "Ôn tập tuần 1: to be, a/an, số nhiều, this/these, chào hỏi và giới thiệu.",
    "activities": [{
        "slug": "checkpoint", "type": "CHECKPOINT", "title": "Kiểm tra cuối tuần 1",
        "instructionsVi": "15 câu ngẫu nhiên từ các bài trong tuần, rồi một cuộc hội thoại với AI (không có gợi ý). "
                          "Cần 70 điểm và hoàn thành nhiệm vụ nói để mở tuần 2.",
        "scenarioSlug": "a1-w1-checkpoint",
        "content": {"fromActivityTypes": ["GRAMMAR", "LISTENING", "READING"], "itemCount": 15, "speaking": True},
        "exercises": [
            fill("We ___ (be) from Vietnam.", ["are"], "We → are.", G1),
            fill("___ she your sister?", ["is"], "Câu hỏi với she → Is.", G1),
            mc("Chọn câu đúng", ["He is an engineer.", "He is a engineer.", "He an engineer."], 0, "an + engineer.", G2),
            fill("I have three ___ (child).", ["children"], "child → children.", G2),
            mc("___ are my colleagues.", ["These", "This", "That"], 0, "Số nhiều → These.", G2),
            error_correction("My father are a doctor.", 2, ["is"], "My father (số ít) → is.", G1),
            transform("They are teachers.", "NEGATIVE", ["They are not teachers.", "They aren't teachers.", "They're not teachers."], "are + not.", G1),
            reorder("Where is the meeting room?", "Where is + the + nơi chốn?", G2),
            translate("Rất vui được gặp bạn.", ["Nice to meet you.", "Nice to meet you!"], "Câu chào lần đầu gặp.", "vocab:greetings"),
            dictation(k + "1", "My name is Huy and I'm a developer.", ""),
            numbers(k + "2", "The meeting room is on the fourteenth floor.", "Phòng họp ở tầng mấy?", "NUMBER", ["14", "14th"], "fourteenth = thứ 14."),
            minimal_pair(k + "1", [("sheep", "/ʃiːp/"), ("ship", "/ʃɪp/")], 1, "'ship' có /ɪ/ ngắn.", "sound:ɪ"),
            ending(k + "1", "watches", "/ɪz/", "Sau 'ch' → /ɪz/."),
        ],
    }],
})

# AI role-play (Tasks 67-70, D10): one scenario per situation lesson.
ROLE_PLAYS = {
    "a1-w1-l1-greetings": "a1-w1-greetings",
    "a1-w1-l2-introducing-yourself": "a1-w1-introducing-yourself",
    "a1-w1-l3-introducing-others": "a1-w1-introducing-others",
    "a1-w1-l4-personal-information": "a1-w1-personal-information",
    "a1-w1-l5-family": "a1-w1-family",
    "a1-w1-l6-hobbies": "a1-w1-hobbies",
    "a1-w1-l7-describing-your-job": "a1-w1-describing-your-job",
    "a1-w1-l8-places-in-the-building": "a1-w1-places-in-the-building",
    "a1-w1-l9-borrowing-things": "a1-w1-borrowing-things",
}

finish(lessons, {
    "slug": "a1-week-1", "title": "Week 1 · Meeting People", "track": "CORE", "orderIndex": 1,
    "description": "Chào hỏi, giới thiệu bản thân và người khác, gia đình, sở thích, nghề nghiệp, công sở. Ngữ pháp: to be, danh từ, a/an, this/these.",
    "icon": "users-round", "emoji": "👋", "color": "indigo",
}, ROLE_PLAYS)
