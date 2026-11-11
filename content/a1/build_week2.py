"""Authoring source for A1 Week 2 · Work and Money.

Grammar: G3 present simple · G4 present continuous · G5 there is/are +
prepositions. Sounds: /eɪ/ /aɪ/ /ɔɪ/ · /aʊ/ /oʊ/ · /p/–/b/ · /t/–/d/ ·
/k/–/ɡ/ · word stress (2 syllables). Helpers: lesson_kit.py.

Usage (from content/a1):
    python3 build_week2.py [--vocab ../a1-vocabulary/out/a1_rows.json]
Outputs: week-2.json, audio_manifest_week-2.tsv
"""
from lesson_kit import *  # noqa: F401,F403 — exercise / activity helpers
from lesson_kit import finish, lesson, start_week

start_week(
    2,
    # Speaker genders per dialogue (A, B, C…), for the TTS voices.
    dialogue_voices={"l1": "FM", "l2": "MF", "l3": "FM", "l4": "FM", "l5": "MF",
                     "l6": "FM", "l7": "FM", "l8": "MF", "l9": "FM"},
    reading_genres={"l1": "PROFILE", "l2": "NOTICE", "l3": "EMAIL", "l4": "NOTICE", "l5": "MENU",
                    "l6": "MESSAGE", "l7": "NOTICE", "l8": "MESSAGE", "l9": "NOTICE"},
)

G3 = "grammar:g3-present-simple"
G4 = "grammar:g4-present-continuous"
G5 = "grammar:g5-there-is-prepositions"
lessons = []

# L1 Daily work schedule --------------------------------------------------------
k = lesson("l1")
lessons.append({
    "slug": "a1-w2-l1-daily-schedule", "title": "Daily work schedule",
    "description": "Nói về một ngày làm việc: giờ giấc, thói quen, always / usually / never.",
    "activities": [
        pronunciation("Âm /eɪ/ và /aɪ/", "/eɪ/ trượt từ 'e' sang 'i' (gần như 'ây'); /aɪ/ trượt từ 'a' sang 'i' (như 'ai'). Phải trượt hết, đừng đọc thành một âm.", ["eɪ", "aɪ"], [
            minimal_pair(k + "1", [("late", "/leɪt/"), ("light", "/laɪt/")], 0, "'late' (muộn) có /eɪ/; 'light' (đèn) có /aɪ/.", "sound:eɪ"),
            minimal_pair(k + "2", [("day", "/deɪ/"), ("die", "/daɪ/")], 1, "'die' có /aɪ/ như 'ai'.", "sound:aɪ"),
            ipa_to_word("/taɪm/", ["time", "team", "tame"], 0, "/aɪ/ thường viết là i_e, y hoặc igh: time, my, night.", "sound:aɪ"),
            word_to_ipa(k + "1", "eight", ["/eɪt/", "/aɪt/", "/et/"], 0, "'eight' có /eɪ/ như trong 'day'.", "sound:eɪ"),
        ]),
        vocabulary("Giờ giấc và thói quen", ["a1-get-up", "a1-breakfast", "a1-start", "a1-finish", "a1-lunch", "a1-get-home",
                                             "a1-early", "a1-late", "a1-usually", "a1-always", "a1-o-clock", "a1-half"], [
            mc("\"I usually get up early\" nghĩa là gì?", ["Tôi thường dậy sớm.", "Tôi luôn dậy muộn.", "Tôi không bao giờ dậy sớm."], 0, "usually = thường; early = sớm.", "vocab:routine"),
            mc("7:30 đọc là gì?", ["half past seven", "half past eight", "seven o'clock"], 0, "half past seven = 7 giờ 30.", "vocab:time"),
            mc("Từ trái nghĩa với 'start'", ["finish", "begin", "get up"], 0, "start (bắt đầu) ↔ finish (kết thúc).", "vocab:routine"),
        ]),
        grammar("Present simple: I work / she works", "Với I / you / we / they dùng động từ nguyên mẫu; với he / she / it thêm -s hoặc -es. Trạng từ tần suất (always, usually, never) đứng trước động từ thường.", "g3-present-simple", [
            fill("I ___ (start) work at eight o'clock.", ["start"], "I + động từ nguyên mẫu.", G3),
            fill("Lan ___ (finish) work at five.", ["finishes"], "Lan = she → finish thêm -es (sau sh).", G3),
            fill("My manager ___ (go) to work by car.", ["goes"], "go → goes với he / she / it.", G3),
            reorder("I usually have lunch at twelve.", "Trạng từ tần suất đứng trước động từ: I usually have…", G3),
            error_correction("He work in a bank.", 1, ["works"], "He + work → works.", G3),
        ]),
        listening(k, "Một ngày làm việc", [
            ("A", "What time do you start work, Nam?"), ("B", "I usually start at eight thirty."),
            ("A", "And what time do you finish?"), ("B", "At half past five. Then I go to the gym."),
            ("A", "Do you have lunch at the office?"), ("B", "Yes, I always have lunch with my team."),
        ], [
            numbers(k + "1", "I usually start at eight thirty.", "Nam bắt đầu làm việc lúc mấy giờ?", "TIME", ["8:30"], "eight thirty = 8:30."),
            true_false("Nam goes to the gym after work.", True, "\"Then I go to the gym.\"", "listening"),
            mc("Nam ăn trưa với ai?", ["his team", "his family", "nobody"], 0, "\"I always have lunch with my team.\"", "listening"),
        ]),
        reading("Một ngày của Hoa", "Hoa is a nurse at Bach Mai Hospital. She works from Monday to Friday. She gets up at six o'clock and starts work at seven. She has lunch at the hospital at half past eleven. She finishes work at four and picks up her son from school. On Saturday she doesn't work.", [
            true_false("Hoa works on Saturday.", False, "\"On Saturday she doesn't work.\"", "reading"),
            mc("Hoa bắt đầu làm việc lúc mấy giờ?", ["7:00", "6:00", "4:00"], 0, "\"…starts work at seven.\"", "reading"),
            mc("Hoa làm nghề gì?", ["nurse", "doctor", "teacher"], 0, "\"Hoa is a nurse.\"", "reading"),
            true_false("Hoa picks up her son after work.", True, "\"She finishes work at four and picks up her son.\"", "reading"),
        ]),
        speaking(["I get up at six thirty.", "I start work at eight o'clock.", "I usually finish work at five."]),
    ],
})

# L2 Asking for leave -----------------------------------------------------------
k = lesson("l2")
lessons.append({
    "slug": "a1-w2-l2-asking-for-leave", "title": "Asking for leave",
    "description": "Xin nghỉ một ngày, nói lý do, sắp xếp người làm thay. Phủ định don't / doesn't.",
    "activities": [
        pronunciation("Âm /aʊ/ và /oʊ/", "/aʊ/ trượt từ 'a' sang 'u' (như 'ao'); /oʊ/ trượt từ 'ô' sang 'u' (như 'âu' nhẹ). 'now' khác 'know'.", ["aʊ", "oʊ"], [
            minimal_pair(k + "1", [("now", "/naʊ/"), ("know", "/noʊ/")], 1, "'know' có /oʊ/ (chữ k câm).", "sound:oʊ"),
            minimal_pair(k + "2", [("town", "/taʊn/"), ("tone", "/toʊn/")], 0, "'town' có /aʊ/ như 'ao'.", "sound:aʊ"),
            ipa_to_word("/haʊs/", ["house", "hose", "horse"], 0, "/aʊ/ thường viết là ou hoặc ow: house, now.", "sound:aʊ"),
            word_to_ipa(k + "1", "home", ["/hoʊm/", "/haʊm/", "/hɑːm/"], 0, "'home' có /oʊ/ như 'go'.", "sound:oʊ"),
        ]),
        vocabulary("Xin nghỉ phép", ["a1-manager", "a1-boss", "a1-holiday", "a1-leave", "a1-tomorrow", "a1-ask",
                                     "a1-busy", "a1-free", "a1-sick", "a1-wednesday", "a1-friday"], [
            mc("\"Can I take a day off?\" dùng để làm gì?", ["Xin nghỉ một ngày", "Xin tăng lương", "Hỏi đường"], 0, "take a day off = nghỉ một ngày.", "vocab:work"),
            mc("\"I'm sick\" nghĩa là gì?", ["Tôi bị ốm", "Tôi bận", "Tôi rảnh"], 0, "sick = ốm.", "vocab:work"),
            mc("Ngày ngay sau Thursday là ngày nào?", ["Friday", "Wednesday", "Saturday"], 0, "Thursday (thứ Năm) → Friday (thứ Sáu).", "vocab:time"),
        ]),
        grammar("Present simple: phủ định don't / doesn't", "I / you / we / they + don't + động từ; he / she / it + doesn't + động từ nguyên mẫu (không thêm -s).", "g3-present-simple", [
            fill("I ___ (not / work) on Saturday.", ["don't work", "do not work"], "I + don't + work.", G3),
            fill("She ___ (not / have) a meeting today.", ["doesn't have", "does not have"], "She + doesn't + have (nguyên mẫu).", G3),
            error_correction("He doesn't works on Friday.", 2, ["work"], "Sau doesn't dùng động từ nguyên mẫu: work.", G3),
            transform("My boss works on Sunday.", "NEGATIVE", ["My boss doesn't work on Sunday.", "My boss does not work on Sunday."], "works → doesn't work.", G3),
            translate("Tôi không làm việc vào thứ Bảy.", ["I don't work on Saturday.", "I do not work on Saturday.", "I don't work on Saturdays."], "don't + work; on + thứ.", G3),
        ]),
        listening(k, "Xin nghỉ một ngày", [
            ("A", "Excuse me, Ms Hoa. Do you have a minute?"), ("B", "Sure. What is it?"),
            ("A", "Can I take a day off on Friday?"), ("B", "On Friday? Why?"),
            ("A", "My sister is getting married."), ("B", "OK, no problem. Who does your work on Friday?"),
            ("A", "Nam does. He knows my project."),
        ], [
            mc("Anh ấy muốn nghỉ ngày nào?", ["Friday", "Thursday", "Monday"], 0, "\"Can I take a day off on Friday?\"", "listening"),
            true_false("The manager says no.", False, "Bà Hoa nói \"OK, no problem.\"", "listening"),
            mc("Ai làm việc thay anh ấy?", ["Nam", "Ms Hoa", "his sister"], 0, "\"Nam does. He knows my project.\"", "listening"),
        ]),
        reading("Quy định nghỉ phép", "Leave policy. Staff have 12 days of paid leave a year. Please ask your manager one week before your day off. If you are sick, call your manager before 9 a.m. You don't need to come to the office. Send a message to HR on the same day.", [
            mc("Mỗi năm nhân viên có bao nhiêu ngày nghỉ có lương?", ["12", "7", "9"], 0, "\"12 days of paid leave a year.\"", "reading"),
            true_false("You ask your manager one day before your day off.", False, "Phải xin trước một tuần (one week before).", "reading"),
            true_false("If you are sick, you call your manager before 9 a.m.", True, "\"call your manager before 9 a.m.\"", "reading"),
            mc("Khi bị ốm, bạn nhắn tin cho ai trong ngày?", ["HR", "your team", "your family"], 0, "\"Send a message to HR on the same day.\"", "reading"),
        ]),
        speaking(["Can I take a day off on Friday?", "My sister is getting married.", "Nam can do my work."]),
    ],
})

# L3 Scheduling a meeting -------------------------------------------------------
k = lesson("l3")
lessons.append({
    "slug": "a1-w2-l3-scheduling-a-meeting", "title": "Scheduling a meeting",
    "description": "Hẹn ngày giờ họp: Are you free…?, on Monday, at two o'clock. Câu hỏi Do / Does.",
    "activities": [
        pronunciation("Âm /p/ và /b/", "/p/ bật hơi mạnh, không rung cổ; /b/ rung cổ, gần như không bật hơi. Đặt tờ giấy trước miệng: /p/ làm giấy rung.", ["p", "b"], [
            minimal_pair(k + "1", [("pack", "/pæk/"), ("back", "/bæk/")], 1, "'back' bắt đầu bằng /b/ rung cổ.", "sound:b"),
            minimal_pair(k + "2", [("pay", "/peɪ/"), ("bay", "/beɪ/")], 0, "'pay' (trả tiền) bật hơi /p/.", "sound:p"),
            minimal_pair(k + "3", [("cap", "/kæp/"), ("cab", "/kæb/")], 0, "Âm cuối /p/ phải bật nhẹ, đừng nuốt mất.", "sound:p"),
            ipa_to_word("/bɪzi/", ["busy", "pizza", "bossy"], 0, "'busy' đọc /bɪzi/ — chữ u đọc thành /ɪ/.", "sound:b"),
        ]),
        vocabulary("Lịch họp", ["a1-meeting", "a1-plan", "a1-project", "a1-team", "a1-monday", "a1-tuesday",
                                "a1-thursday", "a1-morning", "a1-afternoon", "a1-calendar", "a1-minute"], [
            mc("\"Are you free on Tuesday?\" nghĩa là gì?", ["Thứ Ba bạn có rảnh không?", "Thứ Ba bạn có bận không?", "Thứ Ba bạn có đi làm không?"], 0, "free = rảnh.", "vocab:work"),
            mc("Chọn giới từ đúng: ___ Monday", ["on", "in", "at"], 0, "on + thứ trong tuần.", "vocab:time"),
            mc("Chọn giới từ đúng: ___ three o'clock", ["at", "on", "in"], 0, "at + giờ.", "vocab:time"),
        ]),
        grammar("Present simple: câu hỏi Do / Does", "Do + I / you / we / they + động từ…? · Does + he / she / it + động từ nguyên mẫu…? Trả lời ngắn: Yes, I do. / No, she doesn't.", "g3-present-simple", [
            fill("___ you have a meeting on Monday?", ["Do"], "you → Do.", G3),
            fill("___ Peter work on Saturday?", ["Does"], "Peter = he → Does.", G3),
            mc("\"Does she like the plan?\" — trả lời đúng là:", ["Yes, she does.", "Yes, she is.", "Yes, she do."], 0, "Hỏi bằng Does → trả lời does / doesn't.", G3),
            reorder("What time does the meeting start?", "Từ để hỏi + does + chủ ngữ + động từ?", G3),
            transform("They work in Room B.", "QUESTION", ["Do they work in Room B?"], "Thêm Do lên đầu câu.", G3),
        ]),
        listening(k, "Hẹn lịch họp", [
            ("A", "Hi Peter. Are you free on Tuesday afternoon?"), ("B", "Sorry, I'm busy on Tuesday. What about Wednesday?"),
            ("A", "Wednesday is fine. Does two o'clock work for you?"), ("B", "Yes, two o'clock is good. Where is the meeting?"),
            ("A", "In Room B on the third floor."), ("B", "Great. See you on Wednesday."),
        ], [
            mc("Cuộc họp vào ngày nào?", ["Wednesday", "Tuesday", "Thursday"], 0, "Peter bận thứ Ba, họ chọn thứ Tư.", "listening"),
            numbers(k + "1", "Does two o'clock work for you?", "Cuộc họp lúc mấy giờ?", "TIME", ["2:00", "14:00"], "two o'clock = 2 giờ."),
            true_false("The meeting is on the second floor.", False, "Phòng B ở tầng ba (third floor).", "listening"),
        ]),
        reading("Kế hoạch tuần sau", "Hi team, here is our plan for next week. On Monday at 9 a.m. we have the team meeting in Room A. On Wednesday afternoon we meet the customer at their office. On Friday at 4 p.m. we check the project. Please tell me if you are not free. Thanks, Linh", [
            mc("Họp nhóm vào lúc nào?", ["Monday at 9 a.m.", "Friday at 4 p.m.", "Wednesday morning"], 0, "\"On Monday at 9 a.m. we have the team meeting.\"", "reading"),
            true_false("They meet the customer in Room A.", False, "Họ gặp khách ở văn phòng của khách.", "reading"),
            mc("Thứ Sáu cả nhóm làm gì?", ["check the project", "meet the customer", "have a party"], 0, "\"On Friday at 4 p.m. we check the project.\"", "reading"),
            true_false("Linh writes this email.", True, "Email ký tên Linh.", "reading"),
        ]),
        speaking(["Are you free on Tuesday afternoon?", "Let's meet at two o'clock.", "The meeting is in Room B."]),
    ],
})

# L4 Booking a table ------------------------------------------------------------
k = lesson("l4")
lessons.append({
    "slug": "a1-w2-l4-booking-a-table", "title": "Booking a table",
    "description": "Gọi điện đặt bàn: số người, giờ, tên. There is / there are, Is there…?",
    "activities": [
        pronunciation("Âm /t/ và /d/", "/t/ không rung cổ, bật hơi; /d/ rung cổ. Âm cuối /t/, /d/ phải đọc nhẹ nhưng rõ: write ≠ ride.", ["t", "d"], [
            minimal_pair(k + "1", [("two", "/tuː/"), ("do", "/duː/")], 0, "'two' bắt đầu bằng /t/ bật hơi.", "sound:t"),
            minimal_pair(k + "2", [("write", "/raɪt/"), ("ride", "/raɪd/")], 1, "'ride' kết thúc bằng /d/ rung cổ.", "sound:d"),
            minimal_pair(k + "3", [("tie", "/taɪ/"), ("die", "/daɪ/")], 0, "'tie' bắt đầu bằng /t/.", "sound:t"),
            ipa_to_word("/teɪbəl/", ["table", "cable", "label"], 0, "'table' /teɪbəl/: âm đầu /t/, nguyên âm /eɪ/.", "sound:t"),
        ]),
        vocabulary("Đặt bàn", ["a1-restaurant", "a1-table", "a1-people", "a1-window", "a1-tonight", "a1-phone",
                               "a1-evening", "a1-seven", "a1-eight", "a1-name"], [
            mc("\"A table for four\" nghĩa là gì?", ["Bàn cho bốn người", "Bàn số bốn", "Bốn cái bàn"], 0, "a table for + số người.", "vocab:restaurant"),
            mc("\"Tonight\" là khi nào?", ["Tối nay", "Tối mai", "Sáng nay"], 0, "tonight = tối nay.", "vocab:time"),
            mc("Câu nào dùng để đặt bàn?", ["I'd like to book a table.", "I'd like to pay the bill.", "I'd like the menu."], 0, "book a table = đặt bàn.", "vocab:restaurant"),
        ]),
        grammar("There is / there are", "There is + danh từ số ít; there are + danh từ số nhiều. Câu hỏi: Is there…? / Are there…? Phủ định: there isn't / there aren't.", "g5-there-is-prepositions", [
            fill("___ a table for two by the window?", ["Is there"], "Câu hỏi với danh từ số ít: Is there…?", G5),
            fill("There ___ (be) three tables free tonight.", ["are"], "three tables (số nhiều) → are.", G5),
            fill("There ___ (be) a car park behind the restaurant.", ["is"], "a car park (số ít) → is.", G5),
            error_correction("There is two people in my group.", 1, ["are"], "two people (số nhiều) → There are.", G5),
            mc("Chọn câu đúng", ["Are there any tables free?", "Is there any tables free?", "There are any tables free?"], 0, "tables số nhiều → Are there…?", G5),
        ]),
        listening(k, "Đặt bàn qua điện thoại", [
            ("A", "Good evening, Hanoi Kitchen. How can I help you?"), ("B", "Hi. I'd like to book a table for four, please."),
            ("A", "For what time?"), ("B", "Tonight at seven thirty."),
            ("A", "Sorry, there are no tables at seven thirty. Is eight o'clock OK?"), ("B", "Yes, eight is fine. My name is Tuan."),
            ("A", "Thank you, Tuan. See you tonight."),
        ], [
            mc("Anh Tuấn đặt bàn cho mấy người?", ["4", "2", "8"], 0, "\"a table for four\"", "listening"),
            numbers(k + "1", "Is eight o'clock OK?", "Cuối cùng họ đặt bàn lúc mấy giờ?", "TIME", ["8:00", "20:00"], "eight o'clock = 8 giờ (tối)."),
            true_false("There is a free table at seven thirty.", False, "\"there are no tables at seven thirty\"", "listening"),
        ]),
        reading("Thông tin nhà hàng", "Hanoi Kitchen. Open every day from 10 a.m. to 10 p.m. There are 20 tables inside and 6 tables in the garden. There is a big room for parties on the second floor. Please book a table for groups of eight or more. Phone: 024 3826 7788.", [
            mc("Nhà hàng đóng cửa lúc mấy giờ?", ["10 p.m.", "10 a.m.", "8 p.m."], 0, "\"from 10 a.m. to 10 p.m.\"", "reading"),
            true_false("There are six tables in the garden.", True, "\"6 tables in the garden\"", "reading"),
            mc("Phòng tiệc ở đâu?", ["on the second floor", "in the garden", "next to the kitchen"], 0, "\"a big room for parties on the second floor\"", "reading"),
            true_false("A group of four must book a table.", False, "Chỉ nhóm từ 8 người trở lên mới cần đặt.", "reading"),
        ]),
        speaking(["I'd like to book a table for four, please.", "Is there a table by the window?", "My name is Tuan."]),
    ],
})

# L5 Ordering food --------------------------------------------------------------
k = lesson("l5")
lessons.append({
    "slug": "a1-w2-l5-ordering-food", "title": "Ordering food",
    "description": "Gọi món và đồ uống, hỏi món có hay không. like / want / have ở present simple.",
    "activities": [
        pronunciation("Âm /k/ và /ɡ/", "/k/ không rung cổ, bật hơi; /ɡ/ rung cổ. Cả hai: cuống lưỡi chạm vòm họng. Âm cuối: back ≠ bag.", ["k", "ɡ"], [
            minimal_pair(k + "1", [("coat", "/koʊt/"), ("goat", "/ɡoʊt/")], 1, "'goat' bắt đầu bằng /ɡ/ rung cổ.", "sound:ɡ"),
            minimal_pair(k + "2", [("back", "/bæk/"), ("bag", "/bæɡ/")], 1, "Âm cuối /ɡ/ rung cổ, /k/ thì không.", "sound:ɡ"),
            minimal_pair(k + "3", [("cold", "/koʊld/"), ("gold", "/ɡoʊld/")], 0, "'cold' bắt đầu bằng /k/ bật hơi.", "sound:k"),
            ipa_to_word("/tʃɪkɪn/", ["chicken", "kitchen", "checking"], 0, "'chicken' /tʃɪkɪn/ — /k/ ở giữa từ; 'kitchen' là /kɪtʃɪn/.", "sound:k"),
        ]),
        vocabulary("Món ăn và đồ uống", ["a1-menu", "a1-order", "a1-chicken", "a1-beef", "a1-rice", "a1-noodles",
                                         "a1-soup", "a1-vegetable", "a1-water", "a1-coffee", "a1-juice", "a1-delicious"], [
            mc("\"Can I have the menu, please?\" — bạn đang xin gì?", ["thực đơn", "hóa đơn", "nước uống"], 0, "menu = thực đơn.", "vocab:food"),
            mc("Từ nào là đồ uống?", ["juice", "beef", "rice"], 0, "juice = nước ép.", "vocab:food"),
            mc("\"Delicious\" nghĩa là gì?", ["ngon", "đắt", "cay"], 0, "delicious = ngon.", "vocab:food"),
        ]),
        grammar("Present simple: like / want / have", "I like / want / have…; he / she likes / wants / has (have → has). Câu hỏi: Do you want…? Does he like…?", "g3-present-simple", [
            fill("She ___ (like) spicy food.", ["likes"], "She + like → likes.", G3),
            fill("___ you want some water?", ["Do"], "you → Do.", G3),
            fill("My son ___ (not / eat) beef.", ["doesn't eat", "does not eat"], "My son = he → doesn't eat.", G3),
            error_correction("He have rice for lunch.", 1, ["has"], "He + have → has.", G3),
            translate("Tôi muốn một bát phở bò.", ["I want a bowl of beef noodle soup.", "I'd like a bowl of beef noodle soup.", "I want beef noodle soup.", "I'd like beef noodle soup."], "I want / I'd like + món ăn.", G3),
        ]),
        listening(k, "Gọi món", [
            ("A", "Are you ready to order?"), ("B", "Yes. I'd like the chicken with rice, please."),
            ("A", "And to drink?"), ("B", "An orange juice, please. Do you have vegetable soup?"),
            ("A", "Sorry, we don't have soup today."), ("B", "That's OK. Just the chicken and the juice, then."),
        ], [
            mc("Khách gọi món gì?", ["chicken with rice", "beef noodles", "vegetable soup"], 0, "\"I'd like the chicken with rice.\"", "listening"),
            true_false("The restaurant has vegetable soup today.", False, "\"we don't have soup today\"", "listening"),
            mc("Khách uống gì?", ["orange juice", "coffee", "water"], 0, "\"An orange juice, please.\"", "listening"),
        ]),
        reading("Thực đơn bữa trưa", "Lunch menu. Beef noodle soup: 60,000 VND. Chicken with rice: 55,000 VND. Fried vegetables: 35,000 VND. Fresh orange juice: 30,000 VND. Iced coffee: 25,000 VND. Lunch set (chicken with rice and iced coffee): 70,000 VND. We don't have fish on Mondays.", [
            mc("Món nào rẻ nhất?", ["iced coffee", "fried vegetables", "orange juice"], 0, "Cà phê đá 25.000 đồng.", "reading"),
            mc("Phở bò giá bao nhiêu?", ["60,000 VND", "55,000 VND", "70,000 VND"], 0, "\"Beef noodle soup: 60,000 VND\"", "reading"),
            true_false("The lunch set has chicken and coffee.", True, "\"chicken with rice and iced coffee\"", "reading"),
            true_false("They have fish on Mondays.", False, "\"We don't have fish on Mondays.\"", "reading"),
        ]),
        speaking(["Can I have the menu, please?", "I'd like the chicken with rice.", "An orange juice, please."]),
    ],
})

# L6 Paying the bill ------------------------------------------------------------
k = lesson("l6")
lessons.append({
    "slug": "a1-w2-l6-paying-the-bill", "title": "Paying the bill",
    "description": "Xin hóa đơn, hỏi cách thanh toán, tiền thừa. Câu hỏi How much / Where + do / does.",
    "activities": [
        pronunciation("Âm /ɔɪ/ (so với /aɪ/)", "/ɔɪ/ trượt từ 'o' sang 'i' (như 'oi'); /aɪ/ trượt từ 'a' sang 'i' (như 'ai'). boy ≠ buy.", ["ɔɪ", "aɪ"], [
            minimal_pair(k + "1", [("boy", "/bɔɪ/"), ("buy", "/baɪ/")], 1, "'buy' (mua) có /aɪ/.", "sound:aɪ"),
            minimal_pair(k + "2", [("toy", "/tɔɪ/"), ("tie", "/taɪ/")], 0, "'toy' có /ɔɪ/ như 'oi'.", "sound:ɔɪ"),
            ipa_to_word("/kɔɪn/", ["coin", "cone", "can"], 0, "/ɔɪ/ thường viết là oi hoặc oy: coin, boy.", "sound:ɔɪ"),
            word_to_ipa(k + "1", "choice", ["/tʃɔɪs/", "/tʃaɪs/", "/tʃoʊs/"], 0, "'choice' có /ɔɪ/.", "sound:ɔɪ"),
        ]),
        vocabulary("Thanh toán", ["a1-bill", "a1-pay", "a1-cash", "a1-credit-card", "a1-money", "a1-price",
                                  "a1-cost", "a1-change", "a1-hundred", "a1-cafe"], [
            mc("\"Can I have the bill, please?\" nghĩa là gì?", ["Cho tôi xin hóa đơn.", "Cho tôi xem thực đơn.", "Cho tôi thêm nước."], 0, "bill = hóa đơn.", "vocab:money"),
            mc("\"Do you take credit cards?\" nghĩa là gì?", ["Ở đây có nhận thẻ tín dụng không?", "Bạn có tiền mặt không?", "Bạn có thẻ không?"], 0, "take credit cards = nhận thanh toán bằng thẻ.", "vocab:money"),
            mc("\"Keep the change\" nghĩa là gì?", ["Không cần trả lại tiền thừa", "Đổi tiền giúp tôi", "Giữ hóa đơn"], 0, "change ở đây = tiền thừa.", "vocab:money"),
        ]),
        grammar("Câu hỏi Wh- với present simple", "Từ để hỏi + do / does + chủ ngữ + động từ? How much does it cost? Where do I pay?", "g3-present-simple", [
            fill("How much ___ the coffee cost?", ["does"], "the coffee (số ít) → does.", G3),
            fill("Where ___ I pay?", ["do"], "I → do.", G3),
            reorder("How do you want to pay?", "How + do + you + want…?", G3),
            mc("How much ___ these two drinks cost?", ["do", "does", "is"], 0, "these two drinks (số nhiều) → do.", G3),
            error_correction("Where does you pay?", 1, ["do"], "you → do.", G3),
        ]),
        listening(k, "Thanh toán", [
            ("A", "Excuse me, can I have the bill, please?"), ("B", "Sure. That's one hundred and fifty thousand dong."),
            ("A", "Do you take credit cards?"), ("B", "Sorry, only cash today. Our card machine isn't working."),
            ("A", "OK. Here's two hundred thousand."), ("B", "Thank you. Here's your change, fifty thousand."),
        ], [
            numbers(k + "1", "That's one hundred and fifty thousand dong.", "Hóa đơn bao nhiêu nghìn đồng?", "NUMBER", ["150", "150000", "150,000", "150.000"], "one hundred and fifty thousand = 150.000."),
            true_false("The customer pays by credit card.", False, "Hôm nay chỉ nhận tiền mặt.", "listening"),
            mc("Tiền thừa là bao nhiêu?", ["50,000 dong", "15,000 dong", "150,000 dong"], 0, "\"Here's your change, fifty thousand.\"", "listening"),
        ]),
        reading("Hóa đơn quán cà phê", "Coffee House, table 5. Two iced coffees: 50,000. One orange juice: 30,000. One chocolate cake: 45,000. Total: 125,000 VND. Cash: 150,000. Change: 25,000. Thank you and see you again!", [
            mc("Tổng tiền là bao nhiêu?", ["125,000 VND", "150,000 VND", "25,000 VND"], 0, "\"Total: 125,000 VND\"", "reading"),
            mc("Có mấy ly cà phê đá?", ["2", "1", "5"], 0, "\"Two iced coffees\"", "reading"),
            true_false("The customer pays by card.", False, "Khách trả tiền mặt (Cash: 150,000).", "reading"),
            true_false("The cake costs 45,000.", True, "\"One chocolate cake: 45,000\"", "reading"),
        ]),
        speaking(["Can I have the bill, please?", "Do you take credit cards?", "Keep the change."]),
    ],
})

# L7 Asking prices --------------------------------------------------------------
k = lesson("l7")
lessons.append({
    "slug": "a1-w2-l7-asking-prices", "title": "Asking prices",
    "description": "Hỏi giá, mặc cả nhẹ, fifteen hay fifty. Giới từ chỉ vị trí trong cửa hàng.",
    "activities": [
        pronunciation("Trọng âm: fifteen hay fifty?", "Số -teen nhấn âm tiết sau và kéo dài (fif-TEEN); số -ty nhấn âm tiết đầu (FIF-ty). Nghe trọng âm để không nhầm giá tiền.", [], [
            stress(k + "1", "fifteen", ["fif", "teen"], 1, "fif-TEEN: nhấn âm tiết sau."),
            stress(k + "2", "fifty", ["fif", "ty"], 0, "FIF-ty: nhấn âm tiết đầu."),
            minimal_pair(k + "1", [("thirteen", "/θɜːrˈtiːn/"), ("thirty", "/ˈθɜːrti/")], 0, "thir-TEEN nhấn sau và kéo dài /iːn/.", "stress:numbers"),
            minimal_pair(k + "2", [("sixteen", "/sɪkˈstiːn/"), ("sixty", "/ˈsɪksti/")], 1, "SIX-ty nhấn âm tiết đầu.", "stress:numbers"),
        ]),
        vocabulary("Mua sắm và giá cả", ["a1-price", "a1-cheap", "a1-expensive", "a1-buy", "a1-sell", "a1-shop",
                                         "a1-market", "a1-sale", "a1-fifteen", "a1-fifty"], [
            mc("\"It's on sale\" nghĩa là gì?", ["Đang giảm giá", "Hết hàng", "Không bán"], 0, "on sale = đang giảm giá.", "vocab:shopping"),
            mc("Từ trái nghĩa với 'cheap'", ["expensive", "free", "small"], 0, "cheap (rẻ) ↔ expensive (đắt).", "vocab:shopping"),
            mc("Câu hỏi giá nào đúng?", ["How much is this bag?", "How many is this bag?", "What much is this bag?"], 0, "Hỏi giá: How much is…?", "vocab:shopping"),
        ]),
        grammar("Giới từ chỉ vị trí", "in (trong), on (trên), under (dưới), next to (cạnh), behind (sau), between (giữa hai thứ). There is / are + danh từ + giới từ + nơi chốn.", "g5-there-is-prepositions", [
            mc("The bags are ___ the shelf. (trên kệ)", ["on", "in", "under"], 0, "trên bề mặt → on.", G5),
            mc("The shoes are ___ the bags and the hats. (ở giữa)", ["between", "next", "behind"], 0, "ở giữa hai thứ → between.", G5),
            fill("There ___ (be) a cash desk next to the door.", ["is"], "a cash desk (số ít) → is.", G5),
            reorder("The T-shirts are next to the jeans.", "S + are + next to + nơi chốn.", G5),
            translate("Có hai cái áo khoác ở dưới bàn.", ["There are two jackets under the table.", "There are two coats under the table."], "There are + số nhiều + under + the table.", G5),
        ]),
        listening(k, "Hỏi giá ở chợ", [
            ("A", "Excuse me, how much is this bag?"), ("B", "It's fifty thousand dong."),
            ("A", "Fifty or fifteen?"), ("B", "Fifty. Five, zero."),
            ("A", "That's expensive. How much are these hats?"), ("B", "They're on sale. Thirty thousand each."),
            ("A", "OK, I'll take two hats, please."),
        ], [
            numbers(k + "1", "It's fifty thousand dong.", "Cái túi giá bao nhiêu nghìn đồng?", "NUMBER", ["50", "50000", "50,000", "50.000"], "FIF-ty = 50."),
            true_false("The hats are on sale.", True, "\"They're on sale.\"", "listening"),
            mc("Người mua mua gì?", ["two hats", "a bag", "two bags"], 0, "\"I'll take two hats.\"", "listening"),
        ]),
        reading("Giảm giá cuối tuần", "Big sale this weekend! All T-shirts are now 99,000 dong. Jeans are 30% off. There are new bags on the second floor. Shoes are next to the cash desk. Sale times: Saturday and Sunday, 9 a.m. to 9 p.m. Pay by cash or card.", [
            mc("Áo phông giá bao nhiêu?", ["99,000 dong", "30,000 dong", "90,000 dong"], 0, "\"All T-shirts are now 99,000 dong.\"", "reading"),
            true_false("The new bags are on the second floor.", True, "\"There are new bags on the second floor.\"", "reading"),
            true_false("The sale is on Friday.", False, "Giảm giá thứ Bảy và Chủ nhật.", "reading"),
            mc("Giày ở đâu?", ["next to the cash desk", "on the second floor", "behind the bags"], 0, "\"Shoes are next to the cash desk.\"", "reading"),
        ]),
        speaking(["Excuse me, how much is this bag?", "That's too expensive.", "I'll take two, please."]),
    ],
})

# L8 Choosing a size ------------------------------------------------------------
k = lesson("l8")
lessons.append({
    "slug": "a1-w2-l8-choosing-a-size", "title": "Choosing a size",
    "description": "Tìm quần áo, hỏi cỡ và màu, mặc thử. Present continuous: I'm looking for…",
    "activities": [
        pronunciation("Trọng âm từ hai âm tiết", "Danh từ hai âm tiết thường nhấn âm tiết đầu (JACK-et); động từ hai âm tiết thường nhấn âm tiết sau (re-TURN).", [], [
            stress(k + "1", "jacket", ["jack", "et"], 0, "JACK-et: danh từ, nhấn âm tiết đầu."),
            stress(k + "2", "trousers", ["trou", "sers"], 0, "TROU-sers: nhấn âm tiết đầu."),
            stress(k + "3", "return", ["re", "turn"], 1, "re-TURN: động từ, nhấn âm tiết sau."),
            stress(k + "4", "color", ["col", "or"], 0, "COL-or: nhấn âm tiết đầu."),
        ]),
        vocabulary("Quần áo và kích cỡ", ["a1-size", "a1-small", "a1-large", "a1-big", "a1-try-on", "a1-jacket",
                                          "a1-t-shirt", "a1-trousers", "a1-jeans", "a1-color", "a1-wear", "a1-dress-noun"], [
            mc("\"Can I try it on?\" nghĩa là gì?", ["Tôi mặc thử được không?", "Tôi mua được không?", "Tôi đổi được không?"], 0, "try on = mặc thử.", "vocab:clothes"),
            mc("Cỡ nào lớn nhất?", ["large", "medium", "small"], 0, "small < medium < large.", "vocab:clothes"),
            mc("\"It's too small\" nghĩa là gì?", ["Nó quá nhỏ", "Nó vừa", "Nó quá to"], 0, "too + tính từ = quá…", "vocab:clothes"),
        ]),
        grammar("Present continuous: đang làm gì", "am / is / are + V-ing cho việc đang xảy ra lúc nói. I'm looking for a jacket. Are you looking for something?", "g4-present-continuous", [
            fill("I ___ (look) for a blue jacket.", ["am looking", "'m looking"], "I + am + looking.", G4),
            fill("She ___ (try) on a dress.", ["is trying", "'s trying"], "She + is + trying.", G4),
            fill("___ you looking for something?", ["Are"], "Câu hỏi: Are + you + V-ing?", G4),
            error_correction("They is wearing jeans.", 1, ["are"], "They → are.", G4),
            reorder("I am looking for a medium T-shirt.", "S + am + V-ing + …", G4),
        ]),
        listening(k, "Chọn cỡ áo", [
            ("A", "Hello. Are you looking for something?"), ("B", "Yes, I'm looking for a jacket. Do you have this in medium?"),
            ("A", "Yes, here you are. What colour do you want?"), ("B", "Black, please. Can I try it on?"),
            ("A", "Sure. The fitting room is over there."), ("B", "Hmm, it's too small. Do you have a large?"),
        ], [
            mc("Khách đang tìm gì?", ["a jacket", "a T-shirt", "shoes"], 0, "\"I'm looking for a jacket.\"", "listening"),
            mc("Khách muốn màu gì?", ["black", "blue", "white"], 0, "\"Black, please.\"", "listening"),
            true_false("The medium jacket is too big.", False, "Cỡ M quá nhỏ (too small).", "listening"),
        ]),
        reading("Tin nhắn từ cửa hàng", "Hi Mai! I'm at the shop now. I'm looking for a dress for the party. There's a nice red one, but they only have small and large. I usually wear medium. The small is too tight and the large is too big. What do you think? Should I buy the large?", [
            mc("Người viết đang ở đâu?", ["at the shop", "at the party", "at home"], 0, "\"I'm at the shop now.\"", "reading"),
            mc("Người viết thường mặc cỡ nào?", ["medium", "small", "large"], 0, "\"I usually wear medium.\"", "reading"),
            true_false("The shop has the red dress in medium.", False, "Cửa hàng chỉ có cỡ S và L.", "reading"),
            true_false("The small dress is too tight.", True, "\"The small is too tight.\"", "reading"),
        ]),
        speaking(["I'm looking for a jacket.", "Do you have this in medium?", "Can I try it on?"]),
    ],
})

# L9 Returning items ------------------------------------------------------------
k = lesson("l9")
lessons.append({
    "slug": "a1-w2-l9-returning-items", "title": "Returning items",
    "description": "Trả lại hàng lỗi, đổi cỡ, hỏi hóa đơn. Present simple hay present continuous?",
    "activities": [
        pronunciation("Ôn tập âm tuần 2", "Ôn nguyên âm đôi /eɪ/, /aɪ/, /oʊ/ và các cặp /p/–/b/, /t/–/d/, /k/–/ɡ/ ở cuối từ.", ["eɪ", "aɪ", "oʊ"], [
            minimal_pair(k + "1", [("wait", "/weɪt/"), ("white", "/waɪt/")], 0, "'wait' có /eɪ/, 'white' có /aɪ/.", "sound:eɪ"),
            minimal_pair(k + "2", [("cot", "/kɑːt/"), ("got", "/ɡɑːt/")], 1, "'got' bắt đầu bằng /ɡ/.", "sound:ɡ"),
            minimal_pair(k + "3", [("hat", "/hæt/"), ("had", "/hæd/")], 0, "'hat' kết thúc bằng /t/ không rung cổ.", "sound:t"),
            word_to_ipa(k + "1", "broken", ["/broʊkən/", "/braʊkən/", "/brɑːkən/"], 0, "'broken' có /oʊ/ như 'go'.", "sound:oʊ"),
        ]),
        vocabulary("Đổi trả hàng", ["a1-return", "a1-problem", "a1-wrong", "a1-broken", "a1-change", "a1-money",
                                    "a1-keep", "a1-bag", "a1-shop", "a1-buy"], [
            mc("\"I'd like to return this\" nghĩa là gì?", ["Tôi muốn trả lại món này", "Tôi muốn mua thêm", "Tôi muốn mặc thử"], 0, "return = trả lại.", "vocab:shopping"),
            mc("\"It's broken\" nghĩa là gì?", ["Nó bị hỏng", "Nó quá nhỏ", "Nó rất đắt"], 0, "broken = bị hỏng.", "vocab:shopping"),
            mc("\"It doesn't work\" nói về đồ vật nghĩa là gì?", ["Nó không hoạt động", "Nó không vừa", "Nó không đẹp"], 0, "work (máy móc) = chạy, hoạt động.", "vocab:shopping"),
        ]),
        grammar("Present simple hay present continuous?", "Present simple: thói quen, sự thật, quy định (I usually wear M). Present continuous: việc đang xảy ra lúc nói (The phone isn't working now).", "g4-present-continuous", [
            mc("I usually ___ size M.", ["wear", "am wearing", "wearing"], 0, "usually → present simple.", G3),
            mc("Look! The light ___ .", ["is flashing", "flashes", "flash"], 0, "Look! → việc đang xảy ra → present continuous.", G4),
            fill("The phone ___ (not / work) now.", ["isn't working", "is not working"], "now → present continuous.", G4),
            fill("We ___ (not / give) money back after 7 days.", ["don't give", "do not give"], "Quy định chung → present simple.", G3),
            translate("Tôi đang tìm hóa đơn.", ["I'm looking for the receipt.", "I am looking for the receipt.", "I'm looking for my receipt.", "I am looking for my receipt."], "đang + làm → am + V-ing.", G4),
        ]),
        listening(k, "Trả lại hàng", [
            ("A", "Hi. I'd like to return these headphones."), ("B", "What's the problem?"),
            ("A", "The left side isn't working."), ("B", "Do you have the receipt?"),
            ("A", "Yes, here it is. I bought them on Monday."), ("B", "OK. Do you want new ones or your money back?"),
            ("A", "New ones, please."),
        ], [
            mc("Khách trả lại món gì?", ["headphones", "a phone", "a jacket"], 0, "\"I'd like to return these headphones.\"", "listening"),
            true_false("The customer has the receipt.", True, "\"Yes, here it is.\"", "listening"),
            mc("Khách muốn gì?", ["new headphones", "her money back", "a bigger size"], 0, "\"New ones, please.\"", "listening"),
        ]),
        reading("Chính sách đổi trả", "Returns. You can return items within 14 days. Please bring the receipt. Clothes must be new, with the labels on. We don't give money back for sale items, but you can change them for another size or colour. Broken items: please talk to our staff at the cash desk.", [
            mc("Bạn được trả hàng trong bao lâu?", ["14 days", "7 days", "30 days"], 0, "\"within 14 days\"", "reading"),
            true_false("You need the receipt.", True, "\"Please bring the receipt.\"", "reading"),
            true_false("You can get your money back for sale items.", False, "Hàng giảm giá chỉ được đổi cỡ hoặc màu.", "reading"),
            mc("Đồ bị hỏng thì bạn gặp ai?", ["staff at the cash desk", "the manager by email", "nobody"], 0, "\"talk to our staff at the cash desk\"", "reading"),
        ]),
        speaking(["I'd like to return this, please.", "It isn't working.", "Here is the receipt."]),
    ],
})

# Checkpoint lesson (week review) -------------------------------------------------
k = lesson("rv")
lessons.append({
    "slug": "a1-w2-review", "title": "Week 2 review",
    "description": "Ôn tập tuần 2: present simple, present continuous, there is / are, giờ giấc, ăn uống và mua sắm.",
    "activities": [{
        "slug": "checkpoint", "type": "CHECKPOINT", "title": "Kiểm tra cuối tuần 2",
        "instructionsVi": "15 câu ngẫu nhiên từ các bài trong tuần, rồi một cuộc hội thoại với AI (không có gợi ý). "
                          "Cần 70 điểm và hoàn thành nhiệm vụ nói để mở tuần 3.",
        "scenarioSlug": "a1-w2-checkpoint",
        "content": {"fromActivityTypes": ["GRAMMAR", "LISTENING", "READING"], "itemCount": 15, "speaking": True},
        "exercises": [
            fill("She ___ (work) from Monday to Friday.", ["works"], "She + work → works.", G3),
            fill("___ he have a meeting today?", ["Does"], "he → Does.", G3),
            fill("I ___ (not / like) coffee.", ["don't like", "do not like"], "I + don't + like.", G3),
            mc("Look! They ___ lunch now.", ["are having", "have", "has"], 0, "Look! / now → present continuous.", G4),
            fill("___ there a bank near here?", ["Is"], "a bank (số ít) → Is there…?", G5),
            error_correction("There are a cafe next to my office.", 1, ["is"], "a cafe (số ít) → There is.", G5),
            transform("She works on Sunday.", "NEGATIVE", ["She doesn't work on Sunday.", "She does not work on Sunday."], "works → doesn't work.", G3),
            reorder("What time do you finish work?", "What time + do + you + V?", G3),
            translate("Cái áo này bao nhiêu tiền?", ["How much is this shirt?", "How much is this T-shirt?", "How much does this shirt cost?"], "How much is + đồ vật?", "vocab:shopping"),
            dictation(k + "1", "Can I have the bill, please?", ""),
            numbers(k + "2", "The meeting starts at a quarter past nine.", "Cuộc họp bắt đầu lúc mấy giờ?", "TIME", ["9:15"], "a quarter past nine = 9:15."),
            minimal_pair(k + "1", [("pack", "/pæk/"), ("back", "/bæk/")], 0, "'pack' bắt đầu bằng /p/ bật hơi.", "sound:p"),
            stress(k + "1", "fifteen", ["fif", "teen"], 1, "fif-TEEN: nhấn âm tiết sau."),
        ],
    }],
})

# AI role-play (Tasks 67-70, D10): one scenario per situation lesson.
ROLE_PLAYS = {
    "a1-w2-l1-daily-schedule": "a1-w2-daily-schedule",
    "a1-w2-l2-asking-for-leave": "a1-w2-asking-for-leave",
    "a1-w2-l3-scheduling-a-meeting": "a1-w2-scheduling-a-meeting",
    "a1-w2-l4-booking-a-table": "a1-w2-booking-a-table",
    "a1-w2-l5-ordering-food": "a1-w2-ordering-food",
    "a1-w2-l6-paying-the-bill": "a1-w2-paying-the-bill",
    "a1-w2-l7-asking-prices": "a1-w2-asking-prices",
    "a1-w2-l8-choosing-a-size": "a1-w2-choosing-a-size",
    "a1-w2-l9-returning-items": "a1-w2-returning-items",
}

finish(lessons, {
    "slug": "a1-week-2", "title": "Week 2 · Work and Money", "track": "CORE", "orderIndex": 2,
    "description": "Lịch làm việc, xin nghỉ, hẹn họp, nhà hàng, thanh toán, mua sắm và đổi trả. Ngữ pháp: present simple, present continuous, there is / are và giới từ.",
    "icon": "briefcase", "emoji": "💼", "color": "emerald",
}, ROLE_PLAYS)
