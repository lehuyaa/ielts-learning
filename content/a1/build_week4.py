"""Authoring source for A1 Week 4 · Plans and Problems.

Grammar: G8 past simple (regular + common irregular verbs, was / were) ·
G9 future: will / be going to. Sounds: /m/ /n/ /ŋ/ · /w/ · /h/ /j/ ·
past -ed /t/ /d/ /ɪd/ · sentence stress · linking. Helpers: lesson_kit.py.

Usage (from content/a1):
    python3 build_week4.py [--vocab ../a1-vocabulary/out/a1_rows.json]
Outputs: week-4.json, audio_manifest_week-4.tsv
"""
from lesson_kit import *  # noqa: F401,F403 — exercise / activity helpers
from lesson_kit import audio, finish, lesson, start_week

start_week(
    4,
    # Speaker genders per dialogue (A, B, C…), for the TTS voices.
    dialogue_voices={"l1": "FM", "l2": "MF", "l3": "FM", "l4": "FM", "l5": "MF",
                     "l6": "FM", "l7": "FM", "l8": "FM", "l9": "MF"},
    reading_genres={"l1": "MESSAGE", "l2": "MESSAGE", "l3": "MESSAGE", "l4": "EMAIL", "l5": "ARTICLE",
                    "l6": "ARTICLE", "l7": "MESSAGE", "l8": "EMAIL", "l9": "NOTICE"},
)

G7 = "grammar:g7-can-must-should"
G8 = "grammar:g8-past-simple"
G9 = "grammar:g9-future"
ED = "sound:ending-ed"
SENTENCE_STRESS = "stress:sentence"
LINKING = "sound:linking"
lessons = []


def heard(key, sentence, options, correct, explanation, tag):
    """A multiple-choice item played from audio: which sentence / word was it?"""
    return mc("Nghe và chọn câu bạn nghe thấy", options, correct, explanation, tag,
              audio_url=audio(f"{key}-heard", sentence))


def weak_word(key, sentence, options, correct, explanation, question="Từ nào được đọc nhẹ (không nhấn)?"):
    """Sentence stress: which word is weak / stressed in the sentence heard."""
    return mc(f"{question} — \"{sentence}\"", options, correct, explanation, SENTENCE_STRESS,
              audio_url=audio(f"{key}-stress-s", sentence))


# L1 Phone messages -------------------------------------------------------------
k = lesson("l1")
lessons.append({
    "slug": "a1-w4-l1-phone-messages", "title": "Phone messages",
    "description": "Gọi điện, xin gặp ai đó, để lại và ghi lời nhắn. will: I'll tell him, I'll call you back.",
    "activities": [
        pronunciation("Âm /m/ và /n/ ở cuối từ", "/m/: hai môi khép lại; /n/: đầu lưỡi chạm lợi trên, môi mở. Ở cuối từ phải đóng miệng đúng: game ≠ gain, some ≠ sun.", ["m", "n"], [
            minimal_pair(k + "1", [("game", "/ɡeɪm/"), ("gain", "/ɡeɪn/")], 0, "'game' kết thúc bằng /m/: khép môi.", "sound:m"),
            minimal_pair(k + "2", [("mail", "/meɪl/"), ("nail", "/neɪl/")], 0, "'mail' bắt đầu bằng /m/.", "sound:m"),
            minimal_pair(k + "3", [("some", "/sʌm/"), ("sun", "/sʌn/")], 1, "'sun' kết thúc bằng /n/: môi mở, lưỡi chạm lợi.", "sound:n"),
            ipa_to_word("/naɪn/", ["nine", "mine", "line"], 0, "'nine' /naɪn/: /n/ ở đầu và cuối.", "sound:n"),
        ]),
        vocabulary("Gọi điện thoại", ["a1-phone", "a1-message", "a1-call", "a1-mobile-phone", "a1-text", "a1-number",
                                      "a1-later", "a1-answer-verb", "a1-speak", "a1-busy"], [
            mc("\"Can I take a message?\" nghĩa là gì?", ["Tôi ghi lại lời nhắn được không?", "Tôi gửi tin nhắn được không?", "Tôi gọi lại được không?"], 0, "take a message = ghi lại lời nhắn.", "vocab:phone"),
            mc("\"Call me back\" nghĩa là gì?", ["Gọi lại cho tôi", "Nhắn tin cho tôi", "Đừng gọi cho tôi"], 0, "call back = gọi lại.", "vocab:phone"),
            mc("\"Who's calling, please?\" nghĩa là gì?", ["Ai đang gọi vậy ạ?", "Bạn gọi cho ai?", "Bạn gọi lúc mấy giờ?"], 0, "Câu lịch sự để hỏi tên người gọi.", "vocab:phone"),
        ]),
        grammar("will: quyết định ngay và lời hứa", "S + will ('ll) + động từ nguyên mẫu: quyết định ngay lúc nói, lời hứa, đề nghị giúp. Phủ định: won't. I'll tell him. I'll call you back.", "g9-future", [
            fill("Don't worry. I ___ (tell) him.", ["will tell", "'ll tell"], "Lời hứa → will tell.", G9),
            fill("I'm busy now. I ___ call you back later.", ["will", "'ll"], "Hứa gọi lại → will.", G9),
            mc("The phone is ringing. — I ___ answer it.", ["'ll", "am", "did"], 0, "Quyết định ngay lúc nói → I'll.", G9),
            error_correction("She will calls you tomorrow.", 2, ["call"], "Sau will dùng động từ nguyên mẫu: call.", G9),
            translate("Tôi sẽ gọi lại cho bạn.", ["I'll call you back.", "I will call you back.", "I'll call you later.", "I will call you later."], "I'll + call you back.", G9),
        ]),
        listening(k, "Để lại lời nhắn", [
            ("A", "Good morning, ABC Company. Lan speaking."), ("B", "Hello. Can I speak to Mr Nam, please?"),
            ("A", "Sorry, he's in a meeting. Can I take a message?"), ("B", "Yes, please. This is Tom from Sun Bank. Can he call me back?"),
            ("A", "Sure. What's your number?"), ("B", "It's zero nine one two, three four five, six seven eight."),
            ("A", "OK. I'll tell him."),
        ], [
            mc("Ông Nam đang làm gì?", ["in a meeting", "at lunch", "on holiday"], 0, "\"he's in a meeting\"", "listening"),
            numbers(k + "1", "It's zero nine one two, three four five, six seven eight.", "Số điện thoại của Tom là gì?", "PHONE", ["0912345678"], "zero = 0; đọc từng chữ số."),
            true_false("Lan will call Tom back.", False, "Lan sẽ nhắn ông Nam gọi lại cho Tom.", "listening"),
        ]),
        reading("Lời nhắn điện thoại", "Phone message. For: Mr Nam. From: Tom Baker, Sun Bank. Time: 10:15, Tuesday. Message: Tom called about the meeting on Friday. He can't come at 9 a.m. Can you meet at 11 instead? Please call him back today. Number: 0912 345 678. Taken by: Lan", [
            mc("Lời nhắn gửi cho ai?", ["Mr Nam", "Tom Baker", "Lan"], 0, "\"For: Mr Nam.\"", "reading"),
            true_false("Tom can come at 9 a.m. on Friday.", False, "\"He can't come at 9 a.m.\"", "reading"),
            mc("Tom muốn họp lúc mấy giờ?", ["11 a.m.", "9 a.m.", "10:15"], 0, "\"Can you meet at 11 instead?\"", "reading"),
            true_false("Lan took the message.", True, "\"Taken by: Lan\"", "reading"),
        ]),
        speaking(["Can I speak to Mr Nam, please?", "Can I take a message?", "I'll call you back later."]),
    ],
})

# L2 Invitations ----------------------------------------------------------------
k = lesson("l2")
lessons.append({
    "slug": "a1-w4-l2-invitations", "title": "Invitations",
    "description": "Mời và nhận lời mời: Would you like to come?, I'd love to. be going to cho kế hoạch.",
    "activities": [
        pronunciation("Âm /w/ (so với /v/)", "/w/: tròn môi như 'u' rồi mở ra, răng không chạm môi. /v/: răng trên chạm môi dưới. Người Việt hay đọc 'west' thành 'vest'.", ["w", "v"], [
            minimal_pair(k + "1", [("west", "/west/"), ("vest", "/vest/")], 0, "'west' bắt đầu bằng /w/: tròn môi, không chạm răng.", "sound:w"),
            minimal_pair(k + "2", [("wine", "/waɪn/"), ("vine", "/vaɪn/")], 1, "'vine' bắt đầu bằng /v/: răng chạm môi.", "sound:v"),
            minimal_pair(k + "3", [("wet", "/wet/"), ("vet", "/vet/")], 0, "'wet' bắt đầu bằng /w/.", "sound:w"),
            ipa_to_word("/wiːkend/", ["weekend", "weekday", "week"], 0, "'weekend' bắt đầu bằng /w/.", "sound:w"),
        ]),
        vocabulary("Lời mời và sự kiện", ["a1-party", "a1-birthday", "a1-invite", "a1-wedding", "a1-celebrate",
                                          "a1-dinner", "a1-weekend", "a1-free", "a1-fun", "a1-concert"], [
            mc("\"Would you like to come?\" nghĩa là gì?", ["Bạn có muốn đến không?", "Bạn có thích không?", "Bạn đến rồi à?"], 0, "Would you like to + động từ? = lời mời lịch sự.", "vocab:social"),
            mc("\"I'd love to!\" nghĩa là gì?", ["Tôi rất muốn!", "Tôi không thích.", "Tôi đang bận."], 0, "Nhận lời mời: I'd love to!", "vocab:social"),
            mc("\"Invite\" nghĩa là gì?", ["mời", "từ chối", "tổ chức"], 0, "invite = mời.", "vocab:social"),
        ]),
        grammar("be going to: kế hoạch", "am / is / are + going to + động từ nguyên mẫu: kế hoạch đã có từ trước. We're going to have a party on Saturday. Câu hỏi: Are you going to…?", "g9-future", [
            fill("We ___ (have) a party on Saturday. (kế hoạch)", ["are going to have", "'re going to have"], "We + are going to + have.", G9),
            fill("___ you going to come to my birthday?", ["Are"], "you → Are you going to…?", G9),
            mc("She ___ going to invite all her colleagues.", ["is", "are", "will"], 0, "She → is going to.", G9),
            reorder("What are you going to do this weekend?", "Từ để hỏi + are + you + going to + động từ?", G9),
            translate("Tối nay tôi định đi xem hòa nhạc.", ["I'm going to go to a concert tonight.", "I am going to go to a concert tonight.", "I'm going to a concert tonight.", "Tonight I'm going to go to a concert."], "Kế hoạch → I'm going to + go…", G9),
        ]),
        listening(k, "Mời dự tiệc sinh nhật", [
            ("A", "Hi Linh! Are you free on Saturday evening?"), ("B", "I think so. Why?"),
            ("A", "It's my birthday. I'm going to have a party at my house. Would you like to come?"), ("B", "I'd love to! What time?"),
            ("A", "At seven. We're going to have a barbecue in the garden."), ("B", "Great. Can I bring anything?"),
            ("A", "Just bring a drink. See you on Saturday!"),
        ], [
            mc("Bữa tiệc vào lúc nào?", ["Saturday evening", "Sunday evening", "Friday evening"], 0, "\"Are you free on Saturday evening?\"", "listening"),
            numbers(k + "1", "The party starts at seven.", "Tiệc bắt đầu lúc mấy giờ?", "TIME", ["7:00", "19:00"], "seven = 7 giờ (tối)."),
            mc("Linh nên mang gì?", ["a drink", "a cake", "a present"], 0, "\"Just bring a drink.\"", "listening"),
        ]),
        reading("Thiệp mời đám cưới", "You're invited! Mai and Duc are going to get married on Sunday, 20 October. The party is at the Lotus Restaurant, 15 Tran Phu Street, at 6 p.m. There will be dinner, music and dancing. Please tell us if you can come before 10 October. We hope to see you there!", [
            mc("Đám cưới vào ngày nào?", ["Sunday, 20 October", "Saturday, 20 October", "Sunday, 10 October"], 0, "\"on Sunday, 20 October\"", "reading"),
            mc("Tiệc ở đâu?", ["the Lotus Restaurant", "Mai's house", "a hotel"], 0, "\"at the Lotus Restaurant\"", "reading"),
            true_false("You should reply before 10 October.", True, "\"Please tell us… before 10 October.\"", "reading"),
            true_false("There will be no music.", False, "\"There will be dinner, music and dancing.\"", "reading"),
        ]),
        speaking(["Would you like to come to my party?", "I'd love to!", "What time does it start?"]),
    ],
})

# L3 Arranging to meet ----------------------------------------------------------
k = lesson("l3")
lessons.append({
    "slug": "a1-w4-l3-arranging-to-meet", "title": "Arranging to meet",
    "description": "Hẹn gặp: Let's meet…, What about…? will hay going to?",
    "activities": [
        pronunciation("Âm /h/ và /j/", "/h/: thở hơi ra nhẹ ở đầu từ (hat ≠ at). /j/: như 'i' lướt nhanh ở đầu từ (yes, year). Đừng bỏ mất /h/ hay /j/.", ["h", "j"], [
            minimal_pair(k + "1", [("hat", "/hæt/"), ("at", "/æt/")], 0, "'hat' có /h/: thở hơi ra.", "sound:h"),
            minimal_pair(k + "2", [("hill", "/hɪl/"), ("ill", "/ɪl/")], 0, "'hill' bắt đầu bằng /h/.", "sound:h"),
            minimal_pair(k + "3", [("year", "/jɪr/"), ("ear", "/ɪr/")], 0, "'year' bắt đầu bằng /j/.", "sound:j"),
            minimal_pair(k + "4", [("yet", "/jet/"), ("jet", "/dʒet/")], 0, "'yet' có /j/ nhẹ, không phải /dʒ/.", "sound:j"),
        ]),
        vocabulary("Hẹn gặp", ["a1-meet", "a1-plan", "a1-later", "a1-tonight", "a1-tomorrow", "a1-place",
                               "a1-cinema", "a1-cafe", "a1-park", "a1-idea"], [
            mc("\"Let's meet at the cafe\" nghĩa là gì?", ["Mình gặp nhau ở quán cà phê nhé", "Mình đi uống cà phê thôi", "Quán cà phê đóng cửa"], 0, "Let's + động từ = mình cùng… nhé.", "vocab:social"),
            mc("\"What about seven?\" nghĩa là gì?", ["Bảy giờ thì sao?", "Bảy giờ là muộn quá", "Bảy người thì sao?"], 0, "What about…? = … thì sao? (gợi ý).", "vocab:social"),
            mc("\"Good idea!\" nghĩa là gì?", ["Ý hay đấy!", "Chúc may mắn!", "Tôi không biết."], 0, "idea = ý tưởng.", "vocab:social"),
        ]),
        grammar("will hay going to?", "going to: kế hoạch đã quyết từ trước (I booked the tickets). will: quyết định ngay lúc nói, đề nghị, hứa. Let's + động từ nguyên mẫu: rủ cùng làm.", "g9-future", [
            mc("A: 'We have no milk.' B: 'OK, I ___ buy some.' (quyết định ngay)", ["'ll", "'m going to", "did"], 0, "Quyết định ngay lúc nói → will.", G9),
            mc("I ___ visit my parents next weekend. I booked the tickets.", ["'m going to", "'ll", "visited"], 0, "Kế hoạch có sẵn → going to.", G9),
            fill("Let's ___ (meet) at the cafe.", ["meet"], "Let's + động từ nguyên mẫu.", G9),
            fill("Where ___ we meet? — At the station.", ["will", "shall"], "Hỏi tương lai: Where will / shall we…?", G9),
            translate("Chúng ta gặp nhau lúc bảy giờ nhé.", ["Let's meet at seven.", "Let's meet at seven o'clock.", "Let's meet at 7."], "Let's + meet + at seven.", G9),
        ]),
        listening(k, "Hẹn đi xem phim", [
            ("A", "Are you going to see the new film this weekend?"), ("B", "Yes! Let's go together. What about Saturday?"),
            ("A", "Saturday is good. What time?"), ("B", "There's a show at eight. Let's meet at the cafe at seven thirty."),
            ("A", "OK. I'll book the tickets online."), ("B", "Great. See you at the cafe."),
        ], [
            mc("Họ xem phim hôm nào?", ["Saturday", "Sunday", "Friday"], 0, "\"What about Saturday?\" — \"Saturday is good.\"", "listening"),
            numbers(k + "1", "Let's meet at the cafe at seven thirty.", "Họ gặp nhau lúc mấy giờ?", "TIME", ["7:30", "19:30"], "seven thirty = 7:30."),
            true_false("The man will book the tickets.", False, "Người nữ nói \"I'll book the tickets online.\"", "listening"),
        ]),
        reading("Tin nhắn nhóm", "Hi everyone! Our team dinner is going to be on Friday. Let's meet at Pho 24 on Le Loi Street at 6:30. It's next to the bookshop. After dinner, we're going to go to a karaoke bar. Nam will bring his car, so four people can go with him. Please reply by Wednesday. Thanks! Hoa", [
            mc("Bữa tối của nhóm vào hôm nào?", ["Friday", "Wednesday", "Saturday"], 0, "\"Our team dinner is going to be on Friday.\"", "reading"),
            mc("Mọi người gặp nhau lúc mấy giờ?", ["6:30", "6:00", "7:30"], 0, "\"at 6:30\"", "reading"),
            true_false("The restaurant is next to the bookshop.", True, "\"It's next to the bookshop.\"", "reading"),
            true_false("Nam is going to walk.", False, "\"Nam will bring his car.\"", "reading"),
        ]),
        speaking(["Let's meet at the cafe.", "What about Saturday?", "I'll book the tickets."]),
    ],
})

# L4 Changing plans -------------------------------------------------------------
k = lesson("l4")
lessons.append({
    "slug": "a1-w4-l4-changing-plans", "title": "Changing plans",
    "description": "Xin lỗi, đổi lịch hẹn, đề nghị ngày khác. won't, not going to, can't. Đuôi -ed.",
    "activities": [
        pronunciation("Đuôi -ed: /t/, /d/ hay /ɪd/", "Sau /t/, /d/ đọc /ɪd/ (wanted). Sau âm vô thanh /p/ /k/ /s/ /ʃ/ /tʃ/ /f/ đọc /t/ (stopped). Còn lại đọc /d/ (called). Không đọc thành 'ét'.", [], [
            ending(k + "1", "stopped", "/t/", "Sau /p/ vô thanh → /t/: stopt.", ED),
            ending(k + "2", "called", "/d/", "Sau /l/ hữu thanh → /d/: calld.", ED),
            ending(k + "3", "wanted", "/ɪd/", "Sau /t/ → /ɪd/: want-id.", ED),
            ending(k + "4", "changed", "/d/", "Sau /dʒ/ hữu thanh → /d/.", ED),
        ]),
        vocabulary("Đổi kế hoạch", ["a1-change", "a1-instead", "a1-sorry", "a1-problem", "a1-plan", "a1-worry",
                                    "a1-later", "a1-busy", "a1-soon", "a1-wait"], [
            mc("\"I can't make it\" nghĩa là gì?", ["Tôi không đến được", "Tôi không làm được món này", "Tôi không biết"], 0, "can't make it = không đến được (cuộc hẹn).", "vocab:social"),
            mc("\"Instead\" nghĩa là gì?", ["thay vào đó", "cũng vậy", "sau đó"], 0, "Can we meet on Sunday instead? = Gặp Chủ nhật thay vào đó nhé?", "vocab:social"),
            mc("\"Don't worry\" nghĩa là gì?", ["Đừng lo", "Đừng đợi", "Đừng đi"], 0, "worry = lo lắng.", "vocab:social"),
        ]),
        grammar("won't / not going to và can't", "won't + động từ = sẽ không; am / is / are not going to + động từ = không định; can't + động từ = không thể. I'm sorry, I can't come. I won't be late.", "g9-future", [
            fill("Sorry, I ___ come on Friday. I'm sick. (không thể)", ["can't", "cannot"], "không thể → can't.", G7),
            fill("Don't worry. I ___ be late. (sẽ không)", ["won't", "will not"], "Lời hứa phủ định → won't.", G9),
            fill("We ___ (not / go) to the beach. It's raining. (không định)", ["aren't going to go", "are not going to go", "'re not going to go"], "We + aren't going to + go.", G9),
            error_correction("He isn't going to comes.", 4, ["come"], "going to + động từ nguyên mẫu: come.", G9),
            transform("We are going to meet on Monday.", "NEGATIVE", ["We aren't going to meet on Monday.", "We are not going to meet on Monday.", "We're not going to meet on Monday."], "are → aren't.", G9),
        ]),
        listening(k, "Đổi lịch hẹn", [
            ("A", "Hi Tom. I'm really sorry, but I can't meet you on Saturday."), ("B", "Oh no. What's the problem?"),
            ("A", "My manager wants me to work on Saturday."), ("B", "That's OK. Can we meet on Sunday instead?"),
            ("A", "Sunday is perfect. Is ten o'clock OK?"), ("B", "Yes. I'll see you on Sunday at ten."),
        ], [
            mc("Vì sao cô ấy đổi lịch?", ["She has to work.", "She is sick.", "She is on holiday."], 0, "\"My manager wants me to work on Saturday.\"", "listening"),
            mc("Họ hẹn lại vào hôm nào?", ["Sunday", "Saturday", "Monday"], 0, "\"Can we meet on Sunday instead?\"", "listening"),
            numbers(k + "1", "Is ten o'clock OK?", "Hẹn mới lúc mấy giờ?", "TIME", ["10:00"], "ten o'clock = 10 giờ."),
        ]),
        reading("Email đổi lịch họp", "Dear Ms Lan, I'm sorry, but I can't come to our meeting on Thursday at 2 p.m. I have to go to Da Nang for work. Can we meet on Monday instead? I'm free all morning. I will send you the report before the meeting. Sorry for the change. Best regards, Peter", [
            mc("Peter không đến họp được vào hôm nào?", ["Thursday", "Monday", "Friday"], 0, "\"our meeting on Thursday at 2 p.m.\"", "reading"),
            mc("Vì sao Peter không đến được?", ["He has to go to Da Nang.", "He is sick.", "He is on holiday."], 0, "\"I have to go to Da Nang for work.\"", "reading"),
            true_false("Peter is free on Monday afternoon.", False, "Peter rảnh cả buổi sáng thứ Hai (all morning).", "reading"),
            true_false("Peter will send the report before the meeting.", True, "\"I will send you the report before the meeting.\"", "reading"),
        ]),
        speaking(["I'm sorry, I can't come on Saturday.", "Can we meet on Sunday instead?", "Don't worry, I won't be late."]),
    ],
})

# L5 Weather small talk ---------------------------------------------------------
k = lesson("l5")
lessons.append({
    "slug": "a1-w4-l5-weather-small-talk", "title": "Weather small talk",
    "description": "Nói chuyện phiếm về thời tiết hôm nay, hôm qua, ngày mai. was / were và will dự đoán. Nhấn câu.",
    "activities": [
        pronunciation("Trọng âm câu", "Trong câu, từ mang nghĩa (danh từ, động từ chính, tính từ) được nhấn mạnh; từ nhỏ (a, the, is, to, and) đọc nhẹ và nhanh: It's going to RAIN toMORrow.", [], [
            weak_word(k + "1", "It's going to rain tomorrow.", ["to", "rain", "tomorrow"], 0, "'to' là từ nhỏ, đọc nhẹ /tə/; 'rain' và 'tomorrow' được nhấn."),
            weak_word(k + "2", "The weather is nice today.", ["the", "weather", "nice"], 0, "'the' đọc nhẹ /ðə/; 'weather', 'nice', 'today' được nhấn."),
            weak_word(k + "3", "It was cold and windy.", ["cold", "was", "and"], 0, "'cold' và 'windy' mang nghĩa nên được nhấn.", "Từ nào được nhấn mạnh?"),
            stress(k + "1", "weather", ["weath", "er"], 0, "WEATH-er: nhấn âm tiết đầu."),
        ]),
        vocabulary("Thời tiết", ["a1-weather", "a1-sunny", "a1-rainy", "a1-cloudy", "a1-windy", "a1-hot",
                                 "a1-cold", "a1-warm", "a1-cool", "a1-rain", "a1-sun", "a1-snow"], [
            mc("\"It's sunny\" nghĩa là gì?", ["Trời nắng", "Trời mưa", "Trời nhiều mây"], 0, "sunny = có nắng.", "vocab:weather"),
            mc("\"What's the weather like?\" nghĩa là gì?", ["Thời tiết thế nào?", "Bạn thích thời tiết nào?", "Hôm nay là thứ mấy?"], 0, "Câu hỏi thời tiết quen thuộc.", "vocab:weather"),
            mc("Từ trái nghĩa với 'warm'", ["cool", "sunny", "windy"], 0, "warm (ấm) ↔ cool (mát).", "vocab:weather"),
        ]),
        grammar("was / were và will cho dự đoán", "Quá khứ của be: I / he / she / it was; you / we / they were. Dự đoán tương lai: It will rain. It won't be cold.", "g8-past-simple", [
            fill("Yesterday it ___ (be) very hot.", ["was"], "it → was.", G8),
            fill("We ___ (be) at the beach last weekend.", ["were"], "we → were.", G8),
            mc("Tomorrow it ___ rain in the afternoon.", ["will", "was", "is"], 0, "Dự đoán tương lai → will.", G9),
            error_correction("The weather were terrible yesterday.", 2, ["was"], "The weather (số ít) → was.", G8),
            translate("Hôm qua trời lạnh.", ["It was cold yesterday.", "Yesterday it was cold.", "Yesterday was cold."], "Quá khứ của it is → it was.", G8),
        ]),
        listening(k, "Nói chuyện thời tiết", [
            ("A", "Lovely day today, isn't it?"), ("B", "Yes, it's warm and sunny. Yesterday was so cold!"),
            ("A", "I know. And it was very windy too."), ("B", "What's the weather going to be like tomorrow?"),
            ("A", "The app says it will rain in the afternoon."), ("B", "Oh no. I'll take my umbrella, then."),
        ], [
            mc("Hôm nay trời thế nào?", ["warm and sunny", "cold and windy", "rainy"], 0, "\"it's warm and sunny\"", "listening"),
            true_false("It was windy yesterday.", True, "\"And it was very windy too.\"", "listening"),
            mc("Ngày mai trời sẽ thế nào?", ["rain in the afternoon", "snow", "sunny all day"], 0, "\"it will rain in the afternoon\"", "listening"),
        ]),
        reading("Thời tiết cuối tuần", "Weather for Hanoi this weekend. Saturday: sunny and hot, 34°C. It will be a good day for the pool. Sunday: cloudy in the morning, with rain and strong wind in the afternoon, 27°C. Take an umbrella! Next week will be cooler, from 24 to 28°C.", [
            mc("Thứ Bảy nhiệt độ bao nhiêu?", ["34°C", "27°C", "24°C"], 0, "\"Saturday: sunny and hot, 34°C.\"", "reading"),
            true_false("It will rain on Saturday.", False, "Thứ Bảy nắng; Chủ nhật mới mưa.", "reading"),
            mc("Chiều Chủ nhật trời thế nào?", ["rain and strong wind", "sunny and hot", "snow"], 0, "\"with rain and strong wind in the afternoon\"", "reading"),
            true_false("Next week will be warmer.", False, "Tuần sau mát hơn (cooler).", "reading"),
        ]),
        speaking(["Lovely day today, isn't it?", "It was very cold yesterday.", "It's going to rain tomorrow."]),
    ],
})

# L6 Your weekend ---------------------------------------------------------------
k = lesson("l6")
lessons.append({
    "slug": "a1-w4-l6-your-weekend", "title": "Your weekend",
    "description": "Kể về cuối tuần vừa qua: went, had, visited. Past simple có quy tắc và bất quy tắc. Nối âm.",
    "activities": [
        pronunciation("Nối âm", "Phụ âm cuối nối với nguyên âm đầu của từ sau: turn‿it‿off nghe như 'tur-ni-toff', an‿egg như 'a-negg', look‿at như 'loo-kat'. Hiểu nối âm giúp nghe người bản xứ nói nhanh.", [], [
            heard(k + "1", "Turn it off.", ["Turn it off.", "Turn it of.", "Turn off it."], 0, "turn‿it‿off nối thành 'tur-ni-toff'.", LINKING),
            heard(k + "2", "I had an egg.", ["I had an egg.", "I had a egg.", "I have an egg."], 0, "had‿an‿egg nối thành 'ha-da-negg'.", LINKING),
            heard(k + "3", "Look at it.", ["Look at it.", "Look it.", "Looked at."], 0, "look‿at‿it nghe như 'loo-ka-tit'.", LINKING),
            dictation(k + "1", "We went out at eight.", "went‿out‿at‿eight nối liền thành một chuỗi.", LINKING),
        ]),
        vocabulary("Cuối tuần", ["a1-weekend", "a1-visit", "a1-enjoy", "a1-relax", "a1-film", "a1-cook-verb",
                                 "a1-beach", "a1-shopping", "a1-park", "a1-sleep", "a1-go"], [
            mc("\"Did you have a good weekend?\" nghĩa là gì?", ["Cuối tuần của bạn vui không?", "Cuối tuần này bạn làm gì?", "Cuối tuần bạn có rảnh không?"], 0, "Câu hỏi thăm quen thuộc vào thứ Hai.", "vocab:free-time"),
            mc("Quá khứ của 'go' là gì?", ["went", "goed", "gone"], 0, "go → went (bất quy tắc).", G8),
            mc("\"I stayed at home\" nghĩa là gì?", ["Tôi ở nhà", "Tôi về quê", "Tôi đi chơi"], 0, "stay → stayed.", "vocab:free-time"),
        ]),
        grammar("Past simple: worked, went, had", "Động từ có quy tắc: thêm -ed (work → worked, stay → stayed, visit → visited). Bất quy tắc phải học thuộc: go → went, have → had, see → saw, eat → ate, buy → bought, take → took.", "g8-past-simple", [
            fill("Last Saturday I ___ (visit) my grandmother.", ["visited"], "visit → visited.", G8),
            fill("We ___ (go) to the beach on Sunday.", ["went"], "go → went.", G8),
            fill("She ___ (have) dinner with her friends.", ["had"], "have → had.", G8),
            mc("Yesterday I ___ a new phone.", ["bought", "buyed", "buy"], 0, "buy → bought (bất quy tắc).", G8),
            error_correction("I watch a film last night.", 1, ["watched"], "last night → quá khứ: watched.", G8),
        ]),
        listening(k, "Cuối tuần thế nào?", [
            ("A", "Hi Nam! Did you have a good weekend?"), ("B", "Yes, thanks. I went to Ha Long Bay with my family."),
            ("A", "Wow! How was it?"), ("B", "It was beautiful. We took a boat trip and ate seafood."),
            ("A", "Sounds great. I stayed at home and cleaned my apartment."), ("B", "That's OK. Sometimes we need a quiet weekend!"),
        ], [
            mc("Nam đã đi đâu?", ["Ha Long Bay", "Da Nang", "Hue"], 0, "\"I went to Ha Long Bay with my family.\"", "listening"),
            mc("Gia đình Nam đã làm gì?", ["took a boat trip", "went shopping", "played football"], 0, "\"We took a boat trip and ate seafood.\"", "listening"),
            true_false("The woman went out at the weekend.", False, "\"I stayed at home and cleaned my apartment.\"", "listening"),
        ]),
        reading("Cuối tuần của tôi", "My weekend. On Saturday morning I got up late and had breakfast at a cafe near my house. In the afternoon I played badminton with my colleagues. In the evening I watched a film on TV. On Sunday I visited my parents in Bac Ninh. My mother cooked my favourite food. I went home at 8 p.m. It was a relaxing weekend.", [
            mc("Sáng thứ Bảy người viết ăn sáng ở đâu?", ["at a cafe", "at home", "at work"], 0, "\"had breakfast at a cafe near my house\"", "reading"),
            mc("Chiều thứ Bảy người viết làm gì?", ["played badminton", "watched a film", "cooked"], 0, "\"In the afternoon I played badminton.\"", "reading"),
            true_false("On Sunday the writer visited friends in Hanoi.", False, "Người viết về thăm bố mẹ ở Bắc Ninh.", "reading"),
            true_false("The writer went home at 8 p.m.", True, "\"I went home at 8 p.m.\"", "reading"),
        ]),
        speaking(["Did you have a good weekend?", "I went to the beach with my family.", "I stayed at home and relaxed."]),
    ],
})

# L7 Sending a parcel -----------------------------------------------------------
k = lesson("l7")
lessons.append({
    "slug": "a1-w4-l7-sending-a-parcel", "title": "Sending a parcel",
    "description": "Gửi bưu kiện ở bưu điện: nặng bao nhiêu, mất bao lâu, giá. Past simple: Did you…?",
    "activities": [
        pronunciation("Âm cuối /n/ và /ŋ/", "/n/: đầu lưỡi chạm lợi trên; /ŋ/: cuống lưỡi chạm vòm họng, như 'ng'. Đuôi -ing đọc /ɪŋ/, không phải /ɪn/.", ["n", "ŋ"], [
            minimal_pair(k + "1", [("ran", "/ræn/"), ("rang", "/ræŋ/")], 1, "'rang' kết thúc bằng /ŋ/ như 'ng'.", "sound:ŋ"),
            minimal_pair(k + "2", [("win", "/wɪn/"), ("wing", "/wɪŋ/")], 1, "'wing' kết thúc bằng /ŋ/.", "sound:ŋ"),
            minimal_pair(k + "3", [("thin", "/θɪn/"), ("thing", "/θɪŋ/")], 0, "'thin' kết thúc bằng /n/.", "sound:n"),
            ipa_to_word("/sendɪŋ/", ["sending", "send in", "seeing"], 0, "-ing đọc /ɪŋ/.", "sound:ŋ"),
        ]),
        vocabulary("Ở bưu điện", ["a1-post-office", "a1-send", "a1-receive", "a1-mail", "a1-letter", "a1-box",
                                  "a1-address", "a1-heavy", "a1-quickly", "a1-price"], [
            mc("\"Post office\" là gì?", ["bưu điện", "văn phòng", "cửa hàng"], 0, "post office = bưu điện.", "vocab:post"),
            mc("\"How long does it take?\" nghĩa là gì?", ["Mất bao lâu?", "Dài bao nhiêu?", "Bao nhiêu tiền?"], 0, "How long does it take? = mất bao lâu?", "vocab:post"),
            mc("\"It's heavy\" nghĩa là gì?", ["Nó nặng", "Nó to", "Nó đắt"], 0, "heavy = nặng.", "vocab:post"),
        ]),
        grammar("Past simple: câu hỏi Did…?", "Did + chủ ngữ + động từ nguyên mẫu…? Trả lời: Yes, I did. / No, I didn't. Câu hỏi Wh-: When did you send it?", "g8-past-simple", [
            fill("___ you send the parcel yesterday?", ["Did"], "Câu hỏi quá khứ → Did.", G8),
            mc("Did she receive the letter? — Yes, she ___.", ["did", "does", "received"], 0, "Hỏi Did → trả lời did / didn't.", G8),
            fill("When did you ___ (send) it?", ["send"], "Sau did dùng động từ nguyên mẫu.", G8),
            error_correction("Did he sent the box?", 2, ["send"], "Did + động từ nguyên mẫu: send.", G8),
            reorder("Where did you buy the stamps?", "Từ để hỏi + did + chủ ngữ + động từ?", G8),
        ]),
        listening(k, "Gửi bưu kiện", [
            ("A", "Next, please. How can I help?"), ("B", "I'd like to send this parcel to Ho Chi Minh City."),
            ("A", "Please put it on the scale. It's two kilos. Normal or express?"), ("B", "How long does normal take?"),
            ("A", "Three days. Express takes one day, but it's more expensive."), ("B", "Normal is fine. How much is it?"),
            ("A", "Sixty thousand dong."),
        ], [
            numbers(k + "1", "It's two kilos.", "Bưu kiện nặng bao nhiêu ki-lô?", "NUMBER", ["2"], "two kilos = 2 kg."),
            mc("Gửi thường mất bao lâu?", ["three days", "one day", "a week"], 0, "\"Three days.\"", "listening"),
            numbers(k + "2", "Sixty thousand dong.", "Phí gửi bao nhiêu nghìn đồng?", "NUMBER", ["60", "60000", "60,000", "60.000"], "SIX-ty thousand = 60.000."),
        ]),
        reading("Tin nhắn về bưu kiện", "Hi Minh, did you get my parcel? I sent it on Monday from the post office near my office. The man at the post office said it will take three days. Inside there is a book and some coffee from Da Lat. Please tell me when it arrives. Thanks! Lan", [
            mc("Lan gửi bưu kiện hôm nào?", ["Monday", "Thursday", "Sunday"], 0, "\"I sent it on Monday.\"", "reading"),
            mc("Trong bưu kiện có gì?", ["a book and some coffee", "clothes", "a phone"], 0, "\"a book and some coffee from Da Lat\"", "reading"),
            true_false("The parcel will take one day.", False, "\"it will take three days\"", "reading"),
            true_false("Lan sent the parcel from a post office near her office.", True, "\"from the post office near my office\"", "reading"),
        ]),
        speaking(["I'd like to send this parcel.", "How long does it take?", "When did you send it?"]),
    ],
})

# L8 Reporting a problem --------------------------------------------------------
k = lesson("l8")
lessons.append({
    "slug": "a1-w4-l8-reporting-a-problem", "title": "Reporting a problem",
    "description": "Báo đồ hỏng ở khách sạn, nhà, văn phòng: It isn't working, It stopped last night. Past simple: didn't.",
    "activities": [
        pronunciation("Nhấn từ phủ định: can / can't", "Trong câu phủ định, can't / didn't / isn't được nhấn và đọc rõ. 'can' trong câu khẳng định đọc yếu /kən/, còn 'can't' đọc mạnh /kænt/ — nghe trọng âm để biết là được hay không.", [], [
            heard(k + "1", "I can't open the door.", ["I can't open the door.", "I can open the door."], 0, "can't được nhấn và kéo dài /kænt/.", SENTENCE_STRESS),
            heard(k + "2", "I can come at five.", ["I can come at five.", "I can't come at five."], 0, "can đọc yếu /kən/, không nhấn.", SENTENCE_STRESS),
            heard(k + "3", "It didn't work.", ["It didn't work.", "It did work.", "It doesn't work."], 0, "DIDN'T được nhấn mạnh.", SENTENCE_STRESS),
            heard(k + "4", "The light isn't working.", ["The light isn't working.", "The light is working."], 0, "ISN'T được nhấn mạnh.", SENTENCE_STRESS),
        ]),
        vocabulary("Đồ hỏng và sự cố", ["a1-problem", "a1-broken", "a1-help", "a1-air-conditioner", "a1-door",
                                        "a1-key", "a1-machine", "a1-computer", "a1-wrong", "a1-check"], [
            mc("\"The air conditioner doesn't work\" nghĩa là gì?", ["Máy lạnh không chạy", "Máy lạnh quá lạnh", "Máy lạnh mới"], 0, "doesn't work = không hoạt động.", "vocab:problems"),
            mc("\"Can you check it?\" nghĩa là gì?", ["Bạn kiểm tra giúp được không?", "Bạn mua giúp được không?", "Bạn mở giúp được không?"], 0, "check = kiểm tra.", "vocab:problems"),
            mc("\"What's wrong?\" nghĩa là gì?", ["Có chuyện gì vậy?", "Sai ở câu nào?", "Bạn sai rồi."], 0, "What's wrong? = Có vấn đề gì?", "vocab:problems"),
        ]),
        grammar("Past simple: phủ định didn't", "didn't + động từ nguyên mẫu: It didn't work. I didn't break it. Với be: wasn't / weren't.", "g8-past-simple", [
            fill("The printer ___ (not / work) this morning.", ["didn't work", "did not work"], "didn't + work.", G8),
            fill("I ___ (not / break) it. It was old.", ["didn't break", "did not break"], "didn't + break.", G8),
            fill("The water ___ (not / be) hot yesterday.", ["wasn't", "was not"], "be quá khứ phủ định: wasn't.", G8),
            error_correction("It didn't worked after the update.", 2, ["work"], "didn't + động từ nguyên mẫu: work.", G8),
            transform("The light turned on.", "NEGATIVE", ["The light didn't turn on.", "The light did not turn on."], "turned → didn't turn.", G8),
        ]),
        listening(k, "Báo hỏng máy lạnh", [
            ("A", "Reception. How can I help you?"), ("B", "Hi, this is room three oh five. The air conditioner isn't working."),
            ("A", "I'm sorry. When did it stop?"), ("B", "It stopped last night. I turned it off and on, but it didn't help."),
            ("A", "OK. I'll send someone to check it in ten minutes."), ("B", "Thank you. Can you also bring some more towels?"),
        ], [
            numbers(k + "1", "This is room three oh five.", "Khách ở phòng số mấy?", "NUMBER", ["305"], "three oh five = 305 (oh = số 0)."),
            mc("Máy lạnh hỏng từ khi nào?", ["last night", "this morning", "yesterday afternoon"], 0, "\"It stopped last night.\"", "listening"),
            true_false("Turning it off and on fixed the problem.", False, "\"…but it didn't help.\"", "listening"),
        ]),
        reading("Email báo sự cố", "Dear IT team, I am writing about my laptop. Yesterday afternoon it stopped working after an update. I tried to restart it, but the screen is still black. I didn't drop it and it didn't get wet. I have an important meeting at 2 p.m. today. Can you check it this morning, please? Thank you, Minh (Sales, 4th floor)", [
            mc("Laptop hỏng khi nào?", ["yesterday afternoon", "this morning", "last week"], 0, "\"Yesterday afternoon it stopped working.\"", "reading"),
            true_false("Minh dropped the laptop.", False, "\"I didn't drop it.\"", "reading"),
            mc("Minh có cuộc họp lúc mấy giờ?", ["2 p.m.", "10 a.m.", "4 p.m."], 0, "\"an important meeting at 2 p.m. today\"", "reading"),
            true_false("Minh works on the 4th floor.", True, "\"Sales, 4th floor\"", "reading"),
        ]),
        speaking(["The air conditioner isn't working.", "It stopped last night.", "Can you check it, please?"]),
    ],
})

# L9 Emergency help -------------------------------------------------------------
k = lesson("l9")
lessons.append({
    "slug": "a1-w4-l9-emergency-help", "title": "Emergency help",
    "description": "Gọi cấp cứu, nói chuyện gì xảy ra và ở đâu. Ôn past simple và will.",
    "activities": [
        pronunciation("Ôn tập âm tuần 4", "Ôn /ŋ/, /w/, /h/ và đuôi -ed.", ["ŋ", "w", "h"], [
            minimal_pair(k + "1", [("sin", "/sɪn/"), ("sing", "/sɪŋ/")], 1, "'sing' kết thúc bằng /ŋ/.", "sound:ŋ"),
            minimal_pair(k + "2", [("wheel", "/wiːl/"), ("veal", "/viːl/")], 0, "'wheel' bắt đầu bằng /w/ (wh đọc /w/).", "sound:w"),
            minimal_pair(k + "3", [("heat", "/hiːt/"), ("eat", "/iːt/")], 0, "'heat' có /h/ ở đầu.", "sound:h"),
            ending(k + "1", "helped", "/t/", "Sau /p/ vô thanh → /t/: helpt.", ED),
        ]),
        vocabulary("Khẩn cấp", ["a1-help", "a1-police-officer", "a1-hospital", "a1-hurt", "a1-fire",
                                "a1-quickly", "a1-dangerous", "a1-safe", "a1-careful", "a1-happen"], [
            mc("\"What happened?\" nghĩa là gì?", ["Chuyện gì đã xảy ra?", "Bạn đang làm gì?", "Bạn có khỏe không?"], 0, "happen → happened (quá khứ).", "vocab:emergency"),
            mc("\"Call an ambulance!\" nghĩa là gì?", ["Gọi xe cấp cứu!", "Gọi cảnh sát!", "Gọi taxi!"], 0, "ambulance = xe cấp cứu.", "vocab:emergency"),
            mc("Ở Việt Nam, số gọi cấp cứu y tế là số nào?", ["115", "113", "114"], 0, "113 công an, 114 cứu hỏa, 115 cấp cứu.", "vocab:emergency"),
        ]),
        grammar("Ôn tập: past simple và will", "Kể việc vừa xảy ra bằng past simple (He fell. I called…). Phản ứng ngay lúc nói bằng will (I'll call an ambulance).", "g8-past-simple", [
            mc("A man ___ off his motorbike. (vừa xảy ra)", ["fell", "falls", "will fall"], 0, "fall → fell (quá khứ).", G8),
            mc("Don't worry. I ___ call the police.", ["'ll", "did", "was"], 0, "Quyết định ngay → I'll.", G9),
            fill("What ___ (happen)? — There was an accident.", ["happened"], "happen → happened.", G8),
            fill("Where ___ you? — I was at home.", ["were"], "you → were.", G8),
            translate("Tôi sẽ gọi xe cấp cứu.", ["I'll call an ambulance.", "I will call an ambulance."], "I'll + call + an ambulance.", G9),
        ]),
        listening(k, "Gọi cấp cứu", [
            ("A", "Emergency services. What's happening?"), ("B", "A man fell off his motorbike. His leg is hurt."),
            ("A", "Is he awake?"), ("B", "Yes, he is, but he can't stand up."),
            ("A", "Where are you?"), ("B", "On Nguyen Trai Street, near the big supermarket."),
            ("A", "OK. Don't move him. The ambulance will be there in ten minutes."),
        ], [
            mc("Chuyện gì đã xảy ra?", ["A man fell off his motorbike.", "A house is on fire.", "A child is lost."], 0, "\"A man fell off his motorbike.\"", "listening"),
            true_false("The man can stand up.", False, "\"he can't stand up\"", "listening"),
            numbers(k + "1", "The ambulance will be there in ten minutes.", "Bao nhiêu phút nữa xe cấp cứu đến?", "NUMBER", ["10"], "ten minutes = 10 phút."),
        ]),
        reading("Số điện thoại khẩn cấp", "Emergency numbers in Vietnam. Police: 113. Fire: 114. Ambulance: 115. When you call, speak slowly and clearly. Say what happened, where you are and how many people are hurt. Don't put the phone down before the operator tells you. Save these numbers in your phone today.", [
            mc("Gọi cứu hỏa số mấy?", ["114", "113", "115"], 0, "\"Fire: 114.\"", "reading"),
            true_false("You should speak quickly when you call.", False, "\"speak slowly and clearly\"", "reading"),
            mc("Khi gọi cần nói gì?", ["what happened, where you are, how many people are hurt", "only your name", "your phone number"], 0, "\"Say what happened, where you are and how many people are hurt.\"", "reading"),
            true_false("You should save the numbers in your phone.", True, "\"Save these numbers in your phone today.\"", "reading"),
        ]),
        speaking(["Help! There was an accident.", "Please call an ambulance.", "We are near the big supermarket."]),
    ],
})

# Checkpoint lesson (week review) -------------------------------------------------
k = lesson("rv")
lessons.append({
    "slug": "a1-w4-review", "title": "Week 4 review",
    "description": "Ôn tập tuần 4: past simple, will / be going to, gọi điện, hẹn gặp, thời tiết và sự cố.",
    "activities": [{
        "slug": "checkpoint", "type": "CHECKPOINT", "title": "Kiểm tra cuối tuần 4",
        "instructionsVi": "15 câu ngẫu nhiên từ các bài trong tuần, rồi một cuộc hội thoại với AI (không có gợi ý). "
                          "Cần 70 điểm và hoàn thành nhiệm vụ nói để hoàn thành chặng A1.",
        "scenarioSlug": "a1-w4-checkpoint",
        "content": {"fromActivityTypes": ["GRAMMAR", "LISTENING", "READING"], "itemCount": 15, "speaking": True},
        "exercises": [
            fill("Yesterday I ___ (work) until eight.", ["worked"], "work → worked.", G8),
            fill("We ___ (go) to Hue last summer.", ["went"], "go → went.", G8),
            fill("___ you call your mother last night?", ["Did"], "Câu hỏi quá khứ → Did.", G8),
            fill("The meeting ___ (not / start) on time.", ["didn't start", "did not start"], "didn't + start.", G8),
            mc("Look at those clouds! It ___ rain.", ["is going to", "will to", "rains"], 0, "Có dấu hiệu sắp xảy ra → is going to.", G9),
            fill("I'm tired. I think I ___ go to bed early.", ["will", "'ll"], "I think I'll… = quyết định ngay.", G9),
            error_correction("She didn't came to the party.", 2, ["come"], "didn't + động từ nguyên mẫu: come.", G8),
            transform("He visited his parents.", "NEGATIVE", ["He didn't visit his parents.", "He did not visit his parents."], "visited → didn't visit.", G8),
            reorder("What are you going to do tomorrow?", "What + are + you + going to + do?", G9),
            translate("Hôm qua trời mưa.", ["It rained yesterday.", "Yesterday it rained.", "It was rainy yesterday."], "rain → rained.", G8),
            dictation(k + "1", "I'll call you back later.", ""),
            numbers(k + "2", "My number is zero nine eight, seven six five, four three two one.", "Số điện thoại là gì?", "PHONE", ["0987654321"], "Đọc từng chữ số."),
            ending(k + "1", "wanted", "/ɪd/", "Sau /t/ → /ɪd/.", ED),
            minimal_pair(k + "1", [("west", "/west/"), ("vest", "/vest/")], 0, "'west' bắt đầu bằng /w/.", "sound:w"),
        ],
    }],
})

# AI role-play (Tasks 67-70, D10): one scenario per situation lesson.
ROLE_PLAYS = {
    "a1-w4-l1-phone-messages": "a1-w4-phone-messages",
    "a1-w4-l2-invitations": "a1-w4-invitations",
    "a1-w4-l3-arranging-to-meet": "a1-w4-arranging-to-meet",
    "a1-w4-l4-changing-plans": "a1-w4-changing-plans",
    "a1-w4-l5-weather-small-talk": "a1-w4-weather-small-talk",
    "a1-w4-l6-your-weekend": "a1-w4-your-weekend",
    "a1-w4-l7-sending-a-parcel": "a1-w4-sending-a-parcel",
    "a1-w4-l8-reporting-a-problem": "a1-w4-reporting-a-problem",
    "a1-w4-l9-emergency-help": "a1-w4-emergency-help",
}

finish(lessons, {
    "slug": "a1-week-4", "title": "Week 4 · Plans and Problems", "track": "CORE", "orderIndex": 4,
    "description": "Gọi điện, lời mời, hẹn gặp, đổi kế hoạch, thời tiết, cuối tuần, bưu điện, báo sự cố và cấp cứu. Ngữ pháp: past simple, will / be going to.",
    "icon": "calendar-clock", "emoji": "📅", "color": "amber",
}, ROLE_PLAYS)
