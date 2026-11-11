"""Authoring source for A1 Week 3 · Travel and Health.

Grammar: G6 adjectives, comparatives, superlatives · G7 can / must / should.
Sounds: /θ/–/ð/ · /s/–/z/ · /ʃ/–/ʒ/ · /tʃ/–/dʒ/ · /f/–/v/ · /l/–/r/ ·
final consonant clusters · /s/–/ʃ/. Helpers: lesson_kit.py.

Usage (from content/a1):
    python3 build_week3.py [--vocab ../a1-vocabulary/out/a1_rows.json]
Outputs: week-3.json, audio_manifest_week-3.tsv
"""
from lesson_kit import *  # noqa: F401,F403 — exercise / activity helpers
from lesson_kit import finish, lesson, start_week

start_week(
    3,
    # Speaker genders per dialogue (A, B, C…), for the TTS voices.
    dialogue_voices={"l1": "FM", "l2": "MF", "l3": "FM", "l4": "FM", "l5": "FM",
                     "l6": "MF", "l7": "FM", "l8": "FM", "l9": "MF"},
    reading_genres={"l1": "MESSAGE", "l2": "ARTICLE", "l3": "SCHEDULE", "l4": "NOTICE", "l5": "EMAIL",
                    "l6": "SCHEDULE", "l7": "ARTICLE", "l8": "NOTICE", "l9": "NOTICE"},
)

G6 = "grammar:g6-adjectives-comparison"
G7 = "grammar:g7-can-must-should"
CLUSTERS = "sound:clusters"
lessons = []

# L1 Asking for directions ------------------------------------------------------
k = lesson("l1")
lessons.append({
    "slug": "a1-w3-l1-asking-directions", "title": "Asking for directions",
    "description": "Hỏi đường và chỉ đường: turn left, go straight, Can you tell me the way to…?",
    "activities": [
        pronunciation("Âm /θ/ và /ð/", "Đặt đầu lưỡi giữa hai hàm răng rồi thổi hơi. /θ/ không rung cổ (three); /ð/ rung cổ (this). Đừng đọc thành /t/, /s/ hay /d/.", ["θ", "ð"], [
            minimal_pair(k + "1", [("think", "/θɪŋk/"), ("sink", "/sɪŋk/")], 0, "'think' có /θ/: lưỡi giữa hai hàm răng.", "sound:θ"),
            minimal_pair(k + "2", [("three", "/θriː/"), ("tree", "/triː/")], 0, "'three' có /θ/, không phải /t/.", "sound:θ"),
            minimal_pair(k + "3", [("they", "/ðeɪ/"), ("day", "/deɪ/")], 0, "'they' có /ð/ rung cổ, không phải /d/.", "sound:ð"),
            ipa_to_word("/ðɪs/", ["this", "these", "thin"], 0, "'this' /ðɪs/: /ð/ rung cổ, /ɪ/ ngắn.", "sound:ð"),
        ]),
        vocabulary("Chỉ đường", ["a1-left", "a1-right", "a1-straight", "a1-turn", "a1-corner", "a1-street",
                                 "a1-road", "a1-map", "a1-far", "a1-near", "a1-cross", "a1-way"], [
            mc("\"Turn left\" nghĩa là gì?", ["Rẽ trái", "Rẽ phải", "Đi thẳng"], 0, "turn left = rẽ trái.", "vocab:directions"),
            mc("\"Go straight on\" nghĩa là gì?", ["Đi thẳng", "Quay lại", "Rẽ ở góc đường"], 0, "straight on = thẳng tiếp.", "vocab:directions"),
            mc("\"Is it far?\" nghĩa là gì?", ["Có xa không?", "Có gần không?", "Có đắt không?"], 0, "far = xa; near = gần.", "vocab:directions"),
        ]),
        grammar("Can you tell me…? và câu chỉ đường", "Hỏi đường lịch sự: Excuse me, can you tell me the way to…? Chỉ đường bằng câu mệnh lệnh (động từ đứng đầu): Go straight, turn left, take the second street. Sau can luôn là động từ nguyên mẫu.", "g7-can-must-should", [
            fill("___ you tell me the way to the station?", ["Can", "Could"], "Nhờ ai làm gì: Can / Could you…?", G7),
            fill("Go straight and ___ (rẽ) left at the bank.", ["turn"], "turn left = rẽ trái.", G7),
            reorder("Can you show me on the map?", "Can + you + động từ…?", G7),
            error_correction("Can you tells me the way?", 2, ["tell"], "Sau can dùng động từ nguyên mẫu: tell.", G7),
            translate("Bạn có thể nói chậm hơn không?", ["Can you speak more slowly?", "Could you speak more slowly?", "Can you speak more slowly, please?", "Could you speak more slowly, please?", "Can you speak slowly, please?"], "Can / Could you + speak more slowly?", G7),
        ]),
        listening(k, "Hỏi đường đến ga", [
            ("A", "Excuse me, can you tell me the way to the train station?"), ("B", "Sure. Go straight on this street for two hundred metres."),
            ("A", "OK, straight on."), ("B", "Then turn left at the bank. The station is on your right."),
            ("A", "Is it far?"), ("B", "No, it's about five minutes on foot."),
        ], [
            mc("Rẽ trái ở đâu?", ["at the bank", "at the hospital", "at the corner shop"], 0, "\"turn left at the bank\"", "listening"),
            true_false("The station is on the left.", False, "\"The station is on your right.\"", "listening"),
            numbers(k + "1", "No, it's about five minutes on foot.", "Đi bộ mất khoảng bao nhiêu phút?", "NUMBER", ["5"], "five minutes = 5 phút."),
        ]),
        reading("Đường đến văn phòng", "Hi Peter! Welcome to Hanoi. Here's how to get to our office from your hotel. Turn right out of the hotel and walk to the corner. Cross the road at the traffic lights. Take the second street on the left. Our office is the tall glass building next to a big bank. It's on the third floor. Call me if you get lost! Lan", [
            mc("Ra khỏi khách sạn, Peter rẽ hướng nào?", ["right", "left", "straight on"], 0, "\"Turn right out of the hotel.\"", "reading"),
            true_false("Peter crosses the road at the traffic lights.", True, "\"Cross the road at the traffic lights.\"", "reading"),
            mc("Văn phòng ở tầng mấy?", ["third", "second", "first"], 0, "\"It's on the third floor.\"", "reading"),
            true_false("The office is next to a hospital.", False, "Văn phòng cạnh một ngân hàng lớn.", "reading"),
        ]),
        speaking(["Excuse me, where is the train station?", "Can you show me on the map?", "Turn left at the bank."]),
    ],
})

# L2 Taking a taxi --------------------------------------------------------------
k = lesson("l2")
lessons.append({
    "slug": "a1-w3-l2-taking-a-taxi", "title": "Taking a taxi",
    "description": "Nói điểm đến, hỏi thời gian và giá, nhờ tài xế dừng xe. Tính từ, very / too.",
    "activities": [
        pronunciation("Âm /s/ và /z/", "/s/ như 'x' tiếng Việt, không rung cổ; /z/ cùng khẩu hình nhưng rung cổ (như tiếng ong). price ≠ prize.", ["s", "z"], [
            minimal_pair(k + "1", [("sip", "/sɪp/"), ("zip", "/zɪp/")], 1, "'zip' bắt đầu bằng /z/ rung cổ.", "sound:z"),
            minimal_pair(k + "2", [("bus", "/bʌs/"), ("buzz", "/bʌz/")], 0, "'bus' kết thúc bằng /s/ không rung.", "sound:s"),
            minimal_pair(k + "3", [("price", "/praɪs/"), ("prize", "/praɪz/")], 1, "'prize' (giải thưởng) kết thúc bằng /z/.", "sound:z"),
            ipa_to_word("/zuː/", ["zoo", "sue", "so"], 0, "'zoo' /zuː/: âm đầu /z/ rung cổ.", "sound:z"),
        ]),
        vocabulary("Đi taxi", ["a1-taxi", "a1-drive", "a1-car", "a1-airport", "a1-place", "a1-far",
                               "a1-fast", "a1-slow", "a1-safe", "a1-seat", "a1-city"], [
            mc("\"Can you slow down, please?\" nghĩa là gì?", ["Anh đi chậm lại được không?", "Anh đi nhanh lên được không?", "Anh dừng ở đây được không?"], 0, "slow down = chậm lại.", "vocab:travel"),
            mc("\"Can you stop here?\" nghĩa là gì?", ["Dừng ở đây được không?", "Rẽ ở đây được không?", "Đợi ở đây được không?"], 0, "stop = dừng.", "vocab:travel"),
            mc("Từ trái nghĩa với 'fast'", ["slow", "far", "safe"], 0, "fast (nhanh) ↔ slow (chậm).", "vocab:travel"),
        ]),
        grammar("Tính từ: vị trí và very / too", "Tính từ đứng trước danh từ (a fast car) hoặc sau be (The car is fast). Tính từ không thêm -s. very = rất; too = quá (thường là không tốt).", "g6-adjectives-comparison", [
            mc("Chọn câu đúng", ["It's a big city.", "It's a city big.", "It's big a city."], 0, "Tính từ đứng trước danh từ: a big city.", G6),
            fill("The traffic is ___ (quá) bad today. We are late.", ["too"], "too = quá (không tốt).", G6),
            error_correction("The taxis are cheaps here.", 3, ["cheap"], "Tính từ không thêm -s.", G6),
            reorder("This is a very old car.", "a + very + tính từ + danh từ.", G6),
            translate("Tài xế rất thân thiện.", ["The driver is very friendly.", "The taxi driver is very friendly."], "be + very + tính từ.", G6),
        ]),
        listening(k, "Đi taxi ra sân bay", [
            ("A", "Hello. Where to?"), ("B", "The airport, please. How long does it take?"),
            ("A", "About forty minutes. The traffic is very bad now."), ("B", "OK. Can you take the fastest road? My flight is at six."),
            ("A", "Sure. Please put on your seat belt."), ("B", "Thanks. How much is it to the airport?"),
            ("A", "About three hundred thousand dong."),
        ], [
            numbers(k + "1", "About forty minutes.", "Đến sân bay mất khoảng bao nhiêu phút?", "NUMBER", ["40"], "FOR-ty = 40."),
            true_false("The traffic is good now.", False, "\"The traffic is very bad now.\"", "listening"),
            numbers(k + "2", "About three hundred thousand dong.", "Giá khoảng bao nhiêu nghìn đồng?", "NUMBER", ["300", "300000", "300,000", "300.000"], "three hundred thousand = 300.000."),
        ]),
        reading("Đi taxi ở Hà Nội", "Taxis in Hanoi are cheap and easy to find. Use a taxi with a meter, or book a car on an app. App cars are often cheaper. Always tell the driver the address of your hotel. The roads are very busy from 5 to 7 p.m., so a short trip can be slow. Have small notes for the fare.", [
            true_false("Taxis in Hanoi are expensive.", False, "\"Taxis in Hanoi are cheap.\"", "reading"),
            mc("Xe đặt qua app thường thế nào?", ["cheaper", "faster", "safer"], 0, "\"App cars are often cheaper.\"", "reading"),
            mc("Đường đông nhất lúc nào?", ["5 to 7 p.m.", "7 to 9 a.m.", "all day"], 0, "\"The roads are very busy from 5 to 7 p.m.\"", "reading"),
            true_false("It is good to have small notes for the fare.", True, "\"Have small notes for the fare.\"", "reading"),
        ]),
        speaking(["The airport, please.", "How long does it take?", "Can you stop here, please?"]),
    ],
})

# L3 Buying tickets -------------------------------------------------------------
k = lesson("l3")
lessons.append({
    "slug": "a1-w3-l3-buying-tickets", "title": "Buying tickets",
    "description": "Mua vé tàu xe: one-way hay return, giờ đi, giá. So sánh hơn: cheaper than, more expensive than.",
    "activities": [
        pronunciation("Âm /ʃ/ và /ʒ/", "/ʃ/ như khi nói 'suỵt', môi hơi tròn, không rung cổ; /ʒ/ cùng khẩu hình nhưng rung cổ — có trong usually, television.", ["ʃ", "ʒ"], [
            minimal_pair(k + "1", [("sea", "/siː/"), ("she", "/ʃiː/")], 1, "'she' có /ʃ/, môi tròn hơn /s/.", "sound:ʃ"),
            minimal_pair(k + "2", [("sip", "/sɪp/"), ("ship", "/ʃɪp/")], 1, "'ship' bắt đầu bằng /ʃ/.", "sound:ʃ"),
            word_to_ipa(k + "1", "usually", ["/juːʒuəli/", "/juːsuəli/", "/juːʃuəli/"], 0, "'usually' có /ʒ/ rung cổ ở giữa.", "sound:ʒ"),
            ipa_to_word("/telɪvɪʒən/", ["television", "telephone", "telescope"], 0, "'television' kết thúc bằng /ʒən/.", "sound:ʒ"),
        ]),
        vocabulary("Mua vé", ["a1-ticket", "a1-train", "a1-bus", "a1-station", "a1-seat", "a1-return",
                              "a1-travel", "a1-trip", "a1-book", "a1-first", "a1-tourist"], [
            mc("\"A return ticket\" là gì?", ["vé khứ hồi", "vé một chiều", "vé đã trả lại"], 0, "return ticket = vé khứ hồi; one-way = một chiều.", "vocab:travel"),
            mc("\"I'd like to book a seat\" nghĩa là gì?", ["Tôi muốn đặt một chỗ ngồi", "Tôi muốn đổi chỗ", "Tôi muốn ngồi xuống"], 0, "book = đặt trước.", "vocab:travel"),
            mc("\"Trip\" nghĩa là gì?", ["chuyến đi", "vé", "nhà ga"], 0, "trip = chuyến đi.", "vocab:travel"),
        ]),
        grammar("So sánh hơn: cheaper than / more expensive than", "Tính từ ngắn + -er + than (cheap → cheaper than; big → bigger). Tính từ dài: more + tính từ + than (more expensive than). Bất quy tắc: good → better, bad → worse.", "g6-adjectives-comparison", [
            fill("The bus is ___ (cheap) than the train.", ["cheaper"], "cheap → cheaper.", G6),
            fill("The plane is ___ (expensive) than the bus.", ["more expensive"], "Tính từ dài: more expensive.", G6),
            mc("The train is ___ than the bus. (comfortable)", ["more comfortable", "comfortabler", "most comfortable"], 0, "Tính từ dài → more comfortable than.", G6),
            error_correction("The bus is slow than the train.", 3, ["slower"], "slow → slower than.", G6),
            translate("Tàu hỏa tốt hơn xe buýt.", ["The train is better than the bus."], "good → better than.", G6),
        ]),
        listening(k, "Mua vé tàu", [
            ("A", "Hello. Two tickets to Hue, please."), ("B", "One-way or return?"),
            ("A", "Return, please. Is the night train cheaper than the day train?"), ("B", "Yes, it's cheaper, but it's slower. It arrives at eight in the morning."),
            ("A", "OK, the night train, please. How much is it?"), ("B", "It's one million five hundred thousand dong for two return tickets."),
        ], [
            mc("Cô ấy mua mấy vé?", ["2", "1", "4"], 0, "\"Two tickets to Hue.\"", "listening"),
            true_false("The night train is faster.", False, "Tàu đêm rẻ hơn nhưng chậm hơn (slower).", "listening"),
            numbers(k + "1", "It arrives at eight in the morning.", "Tàu đêm đến nơi lúc mấy giờ?", "TIME", ["8:00"], "eight in the morning = 8 giờ sáng."),
        ]),
        reading("Bảng giờ và giá vé", "Hanoi to Hue. Day train: leaves 6:00, arrives 19:30, 450,000 VND. Night train: leaves 22:00, arrives 8:00, 380,000 VND. Bus: leaves 18:00, arrives 6:00, 300,000 VND. Plane: 1 hour 15 minutes, from 1,200,000 VND. Children under 6 travel free on the train.", [
            mc("Cách đi nào rẻ nhất?", ["bus", "night train", "day train"], 0, "Xe khách 300.000 đồng.", "reading"),
            mc("Cách đi nào nhanh nhất?", ["plane", "bus", "day train"], 0, "Máy bay chỉ 1 giờ 15 phút.", "reading"),
            true_false("The night train leaves at 10 p.m.", True, "22:00 = 10 giờ tối.", "reading"),
            true_false("Children under 6 pay half price on the train.", False, "Trẻ dưới 6 tuổi được đi miễn phí.", "reading"),
        ]),
        speaking(["A return ticket to Hue, please.", "Is the bus cheaper than the train?", "What time does it leave?"]),
    ],
})

# L4 Airport check-in -----------------------------------------------------------
k = lesson("l4")
lessons.append({
    "slug": "a1-w3-l4-airport-check-in", "title": "Airport check-in",
    "description": "Làm thủ tục bay: hộ chiếu, hành lý, chỗ ngồi, cổng ra máy bay. must / mustn't.",
    "activities": [
        pronunciation("Âm /tʃ/ và /dʒ/", "/tʃ/ như 'ch' bật hơi, không rung cổ; /dʒ/ rung cổ (như 'giơ' đọc nhanh). cheap ≠ jeep. Âm cuối /dʒ/ trong 'change', 'luggage' phải đọc rõ.", ["tʃ", "dʒ"], [
            minimal_pair(k + "1", [("cheap", "/tʃiːp/"), ("jeep", "/dʒiːp/")], 0, "'cheap' bắt đầu bằng /tʃ/ không rung.", "sound:tʃ"),
            minimal_pair(k + "2", [("chin", "/tʃɪn/"), ("gin", "/dʒɪn/")], 1, "'gin' bắt đầu bằng /dʒ/ rung cổ.", "sound:dʒ"),
            minimal_pair(k + "3", [("rich", "/rɪtʃ/"), ("ridge", "/rɪdʒ/")], 1, "'ridge' kết thúc bằng /dʒ/.", "sound:dʒ"),
            ipa_to_word("/tʃeɪndʒ/", ["change", "chain", "range"], 0, "'change' có /tʃ/ đầu và /ndʒ/ cuối.", "sound:tʃ"),
        ]),
        vocabulary("Ở sân bay", ["a1-airport", "a1-plane", "a1-passport", "a1-check-in", "a1-suitcase", "a1-bag",
                                 "a1-fly", "a1-seat", "a1-window", "a1-foreign"], [
            mc("\"Window or aisle?\" hỏi về gì?", ["chỗ ngồi", "hành lý", "giờ bay"], 0, "window seat = ghế cạnh cửa sổ; aisle seat = ghế cạnh lối đi.", "vocab:travel"),
            mc("\"Passport\" là gì?", ["hộ chiếu", "vé máy bay", "thẻ lên máy bay"], 0, "passport = hộ chiếu.", "vocab:travel"),
            mc("\"Suitcase\" là gì?", ["va li", "túi xách", "ba lô"], 0, "suitcase = va li.", "vocab:travel"),
        ]),
        grammar("must / mustn't: quy định", "must + động từ = phải (quy định, bắt buộc); mustn't + động từ = cấm. must không thêm -s và không có 'to' phía sau.", "g7-can-must-should", [
            fill("You ___ show your passport at check-in. (phải)", ["must"], "bắt buộc → must.", G7),
            fill("You ___ take water through security. (cấm)", ["mustn't", "must not"], "cấm → mustn't.", G7),
            mc("Passengers ___ arrive two hours before the flight.", ["must", "mustn't", "can't"], 0, "Quy định bắt buộc → must.", G7),
            error_correction("She musts check in online.", 1, ["must"], "must không thêm -s.", G7),
            translate("Bạn không được hút thuốc ở đây.", ["You mustn't smoke here.", "You must not smoke here.", "You can't smoke here.", "You cannot smoke here."], "mustn't = cấm.", G7),
        ]),
        listening(k, "Làm thủ tục ở sân bay", [
            ("A", "Good morning. Can I see your passport, please?"), ("B", "Here you are."),
            ("A", "Thank you. Do you have any bags to check in?"), ("B", "Yes, one suitcase."),
            ("A", "Please put it here. Window or aisle seat?"), ("B", "Window, please."),
            ("A", "Here's your boarding pass. Gate twelve. You must be at the gate at nine fifteen."),
        ], [
            numbers(k + "1", "Your gate is gate twelve.", "Cổng ra máy bay số mấy?", "NUMBER", ["12"], "twelve = 12."),
            numbers(k + "2", "You must be at the gate at nine fifteen.", "Phải có mặt ở cổng lúc mấy giờ?", "TIME", ["9:15"], "nine fifteen = 9:15."),
            mc("Hành khách chọn ghế nào?", ["window", "aisle", "middle"], 0, "\"Window, please.\"", "listening"),
        ]),
        reading("Trước khi bay", "Before you fly. You must check in two hours before an international flight. Each passenger can take one bag in the cabin (7 kg). Big suitcases must go in the hold (20 kg). You mustn't take bottles of water through security. You must show your passport and boarding pass at the gate.", [
            mc("Bay quốc tế phải làm thủ tục trước bao lâu?", ["two hours", "one hour", "thirty minutes"], 0, "\"two hours before an international flight\"", "reading"),
            mc("Túi xách tay tối đa bao nhiêu kg?", ["7 kg", "20 kg", "10 kg"], 0, "\"one bag in the cabin (7 kg)\"", "reading"),
            true_false("You can take a bottle of water through security.", False, "\"You mustn't take bottles of water through security.\"", "reading"),
            true_false("You must show your boarding pass at the gate.", True, "\"You must show your passport and boarding pass at the gate.\"", "reading"),
        ]),
        speaking(["Here is my passport.", "I have one suitcase.", "A window seat, please."]),
    ],
})

# L5 Hotel check-in / check-out --------------------------------------------------
k = lesson("l5")
lessons.append({
    "slug": "a1-w3-l5-hotel-check-in", "title": "Hotel check-in and check-out",
    "description": "Nhận và trả phòng, hỏi bữa sáng, nhờ lễ tân. Can I…? / Can you…?",
    "activities": [
        pronunciation("Âm /f/ và /v/", "Răng trên chạm nhẹ môi dưới. /f/ chỉ có hơi, không rung; /v/ rung cổ. Người Việt hay bỏ /v/ cuối từ: five, leave.", ["f", "v"], [
            minimal_pair(k + "1", [("fan", "/fæn/"), ("van", "/væn/")], 1, "'van' bắt đầu bằng /v/ rung cổ.", "sound:v"),
            minimal_pair(k + "2", [("fine", "/faɪn/"), ("vine", "/vaɪn/")], 0, "'fine' bắt đầu bằng /f/ không rung.", "sound:f"),
            minimal_pair(k + "3", [("leaf", "/liːf/"), ("leave", "/liːv/")], 1, "'leave' kết thúc bằng /v/ — đọc rõ, đừng bỏ.", "sound:v"),
            word_to_ipa(k + "1", "five", ["/faɪv/", "/faɪf/", "/vaɪv/"], 0, "'five': /f/ đầu, /v/ cuối.", "sound:v"),
        ]),
        vocabulary("Khách sạn", ["a1-hotel", "a1-room", "a1-night", "a1-book", "a1-stay", "a1-key",
                                 "a1-bed", "a1-floor", "a1-pool", "a1-breakfast", "a1-passport"], [
            mc("\"I have a booking\" nghĩa là gì?", ["Tôi đã đặt phòng", "Tôi muốn đặt phòng", "Tôi có một cuốn sách"], 0, "booking = việc đặt trước.", "vocab:hotel"),
            mc("\"Check out\" ở khách sạn là gì?", ["trả phòng", "nhận phòng", "đặt phòng"], 0, "check in = nhận phòng; check out = trả phòng.", "vocab:hotel"),
            mc("\"Is breakfast included?\" nghĩa là gì?", ["Có bao gồm bữa sáng không?", "Bữa sáng mấy giờ?", "Bữa sáng ở đâu?"], 0, "included = đã bao gồm.", "vocab:hotel"),
        ]),
        grammar("Can I…? / Can you…?", "Can I + động từ…? = xin phép; Can you + động từ…? = nhờ người khác. Lịch sự hơn: Could you…? / May I…? Từ chối: Sorry, you can't.", "g7-can-must-should", [
            fill("___ I have a room with a view?", ["Can", "Could", "May"], "Xin phép: Can / Could / May I…?", G7),
            mc("Bạn muốn nhờ lễ tân gọi taxi. Câu nào đúng?", ["Can you call a taxi for me?", "Can I call a taxi for me?", "Must you call a taxi?"], 0, "Nhờ người khác → Can you…?", G7),
            fill("Sorry, you ___ check in before two o'clock. (không được)", ["can't", "cannot"], "không được → can't.", G7),
            reorder("Can I leave my bag here?", "Can + I + động từ…?", G7),
            translate("Tôi có thể trả phòng muộn không?", ["Can I check out late?", "Could I check out late?", "May I check out late?", "Can I check out late, please?"], "Can I + check out + late?", G7),
        ]),
        listening(k, "Nhận phòng khách sạn", [
            ("A", "Good afternoon. Welcome to the Sea View Hotel."), ("B", "Hi. I have a booking. My name is Minh Tran."),
            ("A", "Yes, a double room for three nights. Can I see your passport?"), ("B", "Here you are. Is breakfast included?"),
            ("A", "Yes. Breakfast is from six to ten on the first floor."), ("B", "Great. Can I check out late on Sunday?"),
            ("A", "Sure, you can check out at one o'clock."),
        ], [
            mc("Anh Minh ở mấy đêm?", ["3", "2", "1"], 0, "\"a double room for three nights\"", "listening"),
            true_false("Breakfast is not included.", False, "\"Is breakfast included?\" — \"Yes.\"", "listening"),
            numbers(k + "1", "Sure, you can check out at one o'clock.", "Anh Minh được trả phòng lúc mấy giờ?", "TIME", ["1:00", "13:00"], "one o'clock = 1 giờ (chiều)."),
        ]),
        reading("Email xác nhận đặt phòng", "Dear Mr Tran, thank you for your booking at the Sea View Hotel. Room: double room with a sea view. Dates: 12 to 15 June (3 nights). Price: 1,500,000 VND a night, breakfast included. Check-in is from 2 p.m. and check-out is before 12 p.m. You can use the pool and the gym for free. We look forward to seeing you!", [
            mc("Giá phòng mỗi đêm là bao nhiêu?", ["1,500,000 VND", "4,500,000 VND", "500,000 VND"], 0, "\"1,500,000 VND a night\"", "reading"),
            true_false("Check-in is from 12 p.m.", False, "Nhận phòng từ 2 giờ chiều.", "reading"),
            true_false("The pool is free.", True, "\"You can use the pool and the gym for free.\"", "reading"),
            mc("Phòng nhìn ra đâu?", ["the sea", "the city", "the garden"], 0, "\"a double room with a sea view\"", "reading"),
        ]),
        speaking(["I have a booking. My name is Minh.", "Is breakfast included?", "Can I check out late?"]),
    ],
})

# L6 Train and bus times --------------------------------------------------------
k = lesson("l6")
lessons.append({
    "slug": "a1-w3-l6-train-bus-times", "title": "Train and bus times",
    "description": "Hỏi giờ tàu xe, chuyến tiếp theo, xuống ở đâu. So sánh nhất: the fastest, the cheapest.",
    "activities": [
        pronunciation("Âm /l/ và /r/", "/l/: đầu lưỡi chạm lợi trên. /r/ (giọng Mỹ): lưỡi cong lên, không chạm đâu. Âm /l/ cuối từ vẫn phải đọc: bill, all.", ["l", "r"], [
            minimal_pair(k + "1", [("light", "/laɪt/"), ("right", "/raɪt/")], 1, "'right' bắt đầu bằng /r/, lưỡi không chạm lợi.", "sound:r"),
            minimal_pair(k + "2", [("late", "/leɪt/"), ("rate", "/reɪt/")], 0, "'late' bắt đầu bằng /l/.", "sound:l"),
            minimal_pair(k + "3", [("long", "/lɔːŋ/"), ("wrong", "/rɔːŋ/")], 1, "'wrong' đọc /rɔːŋ/ (w câm).", "sound:r"),
            ipa_to_word("/treɪn/", ["train", "trail", "tray"], 0, "'train' có cụm /tr/ ở đầu.", "sound:r"),
        ]),
        vocabulary("Giờ tàu xe", ["a1-get-on", "a1-get-off", "a1-miss-verb", "a1-minute", "a1-hour", "a1-early",
                                  "a1-late", "a1-station", "a1-subway", "a1-bus"], [
            mc("\"I missed the bus\" nghĩa là gì?", ["Tôi bị lỡ xe buýt", "Tôi đã lên xe buýt", "Tôi thích xe buýt"], 0, "miss = lỡ (chuyến xe).", "vocab:travel"),
            mc("\"Get off\" là gì?", ["xuống xe", "lên xe", "đổi xe"], 0, "get on = lên xe; get off = xuống xe.", "vocab:travel"),
            mc("\"The next bus\" là gì?", ["chuyến xe buýt tiếp theo", "chuyến xe buýt cuối", "chuyến xe buýt đầu tiên"], 0, "next = tiếp theo.", "vocab:travel"),
        ]),
        grammar("So sánh nhất: the fastest / the most expensive", "Tính từ ngắn: the + -est (fast → the fastest; early → the earliest). Tính từ dài: the most + tính từ (the most expensive). Bất quy tắc: good → the best, bad → the worst.", "g6-adjectives-comparison", [
            fill("The plane is ___ (fast) way to get there.", ["the fastest"], "the + fast + est.", G6),
            fill("This is ___ (expensive) ticket.", ["the most expensive"], "Tính từ dài: the most expensive.", G6),
            mc("Which bus is ___? (early)", ["the earliest", "the earlyest", "the most early"], 0, "early → the earliest (y → i).", G6),
            error_correction("This is the worse bus in the city.", 3, ["worst"], "bad → the worst.", G6),
            translate("Đây là chuyến tàu nhanh nhất.", ["This is the fastest train."], "the fastest + train.", G6),
        ]),
        listening(k, "Hỏi giờ xe buýt", [
            ("A", "Excuse me, which bus goes to the Old Quarter?"), ("B", "Number nine. It's the fastest bus."),
            ("A", "When is the next one?"), ("B", "At ten fifteen. They come every twenty minutes."),
            ("A", "Where do I get off?"), ("B", "Get off at the third stop, near the lake."),
        ], [
            numbers(k + "1", "Take bus number nine.", "Đi xe buýt số mấy?", "NUMBER", ["9"], "number nine = số 9."),
            numbers(k + "2", "The next one is at ten fifteen.", "Chuyến tiếp theo lúc mấy giờ?", "TIME", ["10:15"], "ten fifteen = 10:15."),
            mc("Xuống ở điểm dừng thứ mấy?", ["the third stop", "the first stop", "the last stop"], 0, "\"Get off at the third stop.\"", "listening"),
        ]),
        reading("Giờ chạy tàu điện", "Metro Line 2A: Cat Linh to Yen Nghia. First train: 5:30. Last train: 22:00. Trains come every 10 minutes in the morning and evening, and every 15 minutes at other times. A single ticket costs 8,000 to 15,000 VND. A day ticket costs 30,000 VND. It is the cheapest way to cross the city.", [
            mc("Chuyến cuối lúc mấy giờ?", ["22:00", "5:30", "21:00"], 0, "\"Last train: 22:00.\"", "reading"),
            true_false("In the morning, trains come every 10 minutes.", True, "\"every 10 minutes in the morning and evening\"", "reading"),
            mc("Vé ngày giá bao nhiêu?", ["30,000 VND", "15,000 VND", "8,000 VND"], 0, "\"A day ticket costs 30,000 VND.\"", "reading"),
            true_false("The metro is the most expensive way to cross the city.", False, "Đây là cách rẻ nhất (the cheapest).", "reading"),
        ]),
        speaking(["Which bus goes to the Old Quarter?", "When is the next train?", "Where do I get off?"]),
    ],
})

# L7 Describing symptoms --------------------------------------------------------
k = lesson("l7")
lessons.append({
    "slug": "a1-w3-l7-describing-symptoms", "title": "Describing symptoms",
    "description": "Nói mình bị đau ở đâu, đặt lịch khám. Lời khuyên với should / shouldn't.",
    "activities": [
        pronunciation("Cụm phụ âm cuối", "Tiếng Anh có nhiều cụm phụ âm cuối: cold /ld/, hands /ndz/, next /kst/, asked /skt/. Đọc đủ từng âm, đừng bỏ âm cuối: cold ≠ coal.", [], [
            minimal_pair(k + "1", [("cold", "/koʊld/"), ("coal", "/koʊl/")], 0, "'cold' kết thúc bằng /ld/ — có /d/.", CLUSTERS),
            minimal_pair(k + "2", [("hand", "/hænd/"), ("hands", "/hændz/")], 1, "'hands' kết thúc bằng /ndz/.", CLUSTERS),
            minimal_pair(k + "3", [("ask", "/æsk/"), ("asked", "/æskt/")], 1, "'asked' kết thúc bằng /skt/.", CLUSTERS),
            word_to_ipa(k + "1", "next", ["/nekst/", "/neks/", "/net/"], 0, "'next' kết thúc bằng /kst/.", CLUSTERS),
        ]),
        vocabulary("Triệu chứng", ["a1-headache", "a1-fever", "a1-cough", "a1-cold-noun", "a1-pain", "a1-hurt",
                                   "a1-sick", "a1-doctor", "a1-appointment", "a1-rest", "a1-stomach", "a1-tooth"], [
            mc("\"I have a headache\" nghĩa là gì?", ["Tôi bị đau đầu", "Tôi bị đau bụng", "Tôi bị sốt"], 0, "headache = đau đầu.", "vocab:health"),
            mc("\"My back hurts\" nghĩa là gì?", ["Tôi bị đau lưng", "Tôi bị đau tay", "Tôi bị ho"], 0, "hurt = đau.", "vocab:health"),
            mc("\"Fever\" là gì?", ["sốt", "ho", "cảm lạnh"], 0, "fever = sốt; cough = ho; a cold = cảm lạnh.", "vocab:health"),
        ]),
        grammar("should / shouldn't: lời khuyên", "should + động từ = nên; shouldn't + động từ = không nên. Hỏi ý kiến: Should I…? should không thêm -s.", "g7-can-must-should", [
            fill("You have a fever. You ___ see a doctor.", ["should"], "Khuyên nên làm → should.", G7),
            fill("You ___ go to work today. Stay at home. (không nên)", ["shouldn't", "should not"], "không nên → shouldn't.", G7),
            mc("I have a toothache. What ___ I do?", ["should", "must", "am"], 0, "Hỏi ý kiến: What should I do?", G7),
            error_correction("She shoulds rest at home.", 1, ["should"], "should không thêm -s.", G7),
            translate("Bạn nên nghỉ ngơi.", ["You should rest.", "You should have a rest.", "You should get some rest."], "should + rest.", G7),
        ]),
        listening(k, "Gọi điện cho phòng khám", [
            ("A", "City Clinic. How can I help you?"), ("B", "Hi. I feel sick. I have a fever and a bad cough."),
            ("A", "Do you have a headache too?"), ("B", "Yes, and my throat hurts."),
            ("A", "You should see the doctor. Can you come at four o'clock today?"), ("B", "Yes, four o'clock is fine."),
            ("A", "OK. You should drink a lot of water and rest."),
        ], [
            mc("Anh ấy bị gì?", ["fever and cough", "stomach ache", "a broken arm"], 0, "\"I have a fever and a bad cough.\"", "listening"),
            numbers(k + "1", "Can you come at four o'clock today?", "Lịch hẹn lúc mấy giờ?", "TIME", ["4:00", "16:00"], "four o'clock = 4 giờ."),
            true_false("The receptionist says he should go to work.", False, "Cô ấy khuyên uống nhiều nước và nghỉ ngơi.", "listening"),
        ]),
        reading("Bị cảm thì làm gì?", "Do you have a cold? Here is some advice. You should rest and sleep a lot. You should drink warm water or tea with honey. You shouldn't go to the gym. Wash your hands often, and don't share cups with your family. If you have a high fever for more than three days, you should see a doctor.", [
            true_false("You should go to the gym.", False, "\"You shouldn't go to the gym.\"", "reading"),
            mc("Nên uống gì?", ["warm water or tea", "cold juice", "coffee"], 0, "\"drink warm water or tea with honey\"", "reading"),
            mc("Khi nào nên đi khám?", ["high fever for more than three days", "after one day", "never"], 0, "\"If you have a high fever for more than three days…\"", "reading"),
            true_false("You shouldn't share cups with your family.", True, "\"don't share cups with your family\"", "reading"),
        ]),
        speaking(["I have a headache and a fever.", "My back hurts.", "Should I see a doctor?"]),
    ],
})

# L8 Buying medicine ------------------------------------------------------------
k = lesson("l8")
lessons.append({
    "slug": "a1-w3-l8-buying-medicine", "title": "Buying medicine",
    "description": "Mua thuốc, hỏi cách dùng: uống mấy viên, mấy lần. Ôn can / must / should.",
    "activities": [
        pronunciation("Âm /s/ và /ʃ/", "/s/ như 'x'; /ʃ/ như khi nói 'suỵt', môi tròn và đẩy ra trước. Người Việt hay đọc 'she' thành 'see'.", ["s", "ʃ"], [
            minimal_pair(k + "1", [("see", "/siː/"), ("she", "/ʃiː/")], 1, "'she' có /ʃ/, môi tròn.", "sound:ʃ"),
            minimal_pair(k + "2", [("save", "/seɪv/"), ("shave", "/ʃeɪv/")], 0, "'save' bắt đầu bằng /s/.", "sound:s"),
            minimal_pair(k + "3", [("sort", "/sɔːrt/"), ("short", "/ʃɔːrt/")], 1, "'short' bắt đầu bằng /ʃ/.", "sound:ʃ"),
            ipa_to_word("/fɑːrməsi/", ["pharmacy", "farmer", "family"], 0, "'pharmacy': ph đọc là /f/.", "sound:f"),
        ]),
        vocabulary("Ở hiệu thuốc", ["a1-medicine", "a1-pharmacy", "a1-take", "a1-careful", "a1-health", "a1-healthy",
                                    "a1-ill", "a1-well", "a1-help-noun", "a1-pain"], [
            mc("\"Take two tablets a day\" nghĩa là gì?", ["Uống hai viên mỗi ngày", "Mua hai hộp mỗi ngày", "Uống hai lần mỗi tuần"], 0, "take (thuốc) = uống.", "vocab:health"),
            mc("\"Pharmacy\" là gì?", ["hiệu thuốc", "bệnh viện", "phòng khám"], 0, "pharmacy = hiệu thuốc.", "vocab:health"),
            mc("\"Before meals\" nghĩa là gì?", ["trước bữa ăn", "sau bữa ăn", "trong bữa ăn"], 0, "before = trước; after = sau.", "vocab:health"),
        ]),
        grammar("can / must / should: ôn tập", "can = có thể, được phép; must = phải (bắt buộc); mustn't = cấm; should = nên; shouldn't = không nên. Sau cả ba luôn là động từ nguyên mẫu.", "g7-can-must-should", [
            mc("You ___ take this medicine with food. (bác sĩ bắt buộc)", ["must", "can", "should"], 0, "Bắt buộc → must.", G7),
            mc("You ___ drive after taking this medicine. (cấm)", ["mustn't", "shouldn't", "can"], 0, "Cấm → mustn't.", G7),
            fill("___ I take it with milk?", ["Can", "Should", "May"], "Hỏi xin phép / hỏi ý kiến: Can / Should I…?", G7),
            fill("Children ___ take this medicine. It's only for adults. (không được)", ["can't", "cannot", "mustn't", "must not"], "không được → can't / mustn't.", G7),
            reorder("You should take it after meals.", "should + động từ nguyên mẫu.", G7),
        ]),
        listening(k, "Ở hiệu thuốc", [
            ("A", "Hello. Can I help you?"), ("B", "Yes. My throat hurts. Can I have something for it?"),
            ("A", "Try these. Take one tablet three times a day after meals."), ("B", "Can I take them with coffee?"),
            ("A", "No, you shouldn't. Take them with water."), ("B", "OK. How much are they?"),
            ("A", "Forty-five thousand dong."),
        ], [
            numbers(k + "1", "Take one tablet three times a day after meals.", "Uống mấy lần một ngày?", "NUMBER", ["3"], "three times a day = 3 lần một ngày."),
            true_false("He can take the tablets with coffee.", False, "\"No, you shouldn't. Take them with water.\"", "listening"),
            numbers(k + "2", "Forty-five thousand dong.", "Thuốc giá bao nhiêu nghìn đồng?", "NUMBER", ["45", "45000", "45,000", "45.000"], "forty-five thousand = 45.000."),
        ]),
        reading("Hướng dẫn trên hộp thuốc", "Cold and flu tablets. Adults: take 2 tablets every 6 hours. Do not take more than 8 tablets in 24 hours. Children under 12 must not take this medicine. Take with water after food. You should not drink alcohol with these tablets. Keep in a cool, dry place.", [
            mc("Người lớn uống mấy viên mỗi lần?", ["2", "6", "8"], 0, "\"take 2 tablets every 6 hours\"", "reading"),
            true_false("Children under 12 can take this medicine.", False, "\"Children under 12 must not take this medicine.\"", "reading"),
            mc("Tối đa bao nhiêu viên trong 24 giờ?", ["8", "6", "12"], 0, "\"Do not take more than 8 tablets in 24 hours.\"", "reading"),
            true_false("You should take the tablets after food.", True, "\"Take with water after food.\"", "reading"),
        ]),
        speaking(["My throat hurts.", "Can I have something for it?", "How often should I take it?"]),
    ],
})

# L9 At the bank ----------------------------------------------------------------
k = lesson("l9")
lessons.append({
    "slug": "a1-w3-l9-at-the-bank", "title": "At the bank",
    "description": "Đổi tiền, hỏi tỷ giá, dùng ATM. Ôn so sánh và can / must.",
    "activities": [
        pronunciation("Ôn tập âm tuần 3", "Ôn /θ/, /v/, /r/, /dʒ/ và âm cuối: lưỡi giữa răng, môi chạm răng, lưỡi cong, đọc đủ âm cuối.", ["θ", "v", "r"], [
            minimal_pair(k + "1", [("three", "/θriː/"), ("free", "/friː/")], 0, "'three' có /θ/ (lưỡi giữa răng), 'free' có /f/.", "sound:θ"),
            minimal_pair(k + "2", [("very", "/veri/"), ("berry", "/beri/")], 0, "'very' bắt đầu bằng /v/: răng chạm môi dưới.", "sound:v"),
            minimal_pair(k + "3", [("collect", "/kəlekt/"), ("correct", "/kərekt/")], 1, "'correct' có /r/ ở giữa.", "sound:r"),
            word_to_ipa(k + "1", "change", ["/tʃeɪndʒ/", "/ʃeɪndʒ/", "/tʃeɪnʒ/"], 0, "/tʃ/ ở đầu, /ndʒ/ ở cuối.", "sound:dʒ"),
        ]),
        vocabulary("Ở ngân hàng", ["a1-bank", "a1-money", "a1-cash", "a1-change", "a1-form", "a1-sign-verb",
                                   "a1-password", "a1-code", "a1-dollar", "a1-number"], [
            mc("\"Can I change dollars here?\" nghĩa là gì?", ["Tôi đổi đô la ở đây được không?", "Tôi gửi đô la được không?", "Đô la ở đây giá bao nhiêu?"], 0, "change money = đổi tiền.", "vocab:money"),
            mc("\"Please sign here\" nghĩa là gì?", ["Vui lòng ký vào đây", "Vui lòng viết tên", "Vui lòng đợi ở đây"], 0, "sign = ký tên.", "vocab:money"),
            mc("\"Password\" là gì?", ["mật khẩu", "số tài khoản", "chữ ký"], 0, "password = mật khẩu.", "vocab:money"),
        ]),
        grammar("Ôn tập: can / must và so sánh", "must = bắt buộc, mustn't = cấm, can = được phép. So sánh hơn: -er than / more … than; bất quy tắc: good → better.", "g7-can-must-should", [
            mc("You ___ show your ID card to open an account.", ["must", "mustn't", "shouldn't"], 0, "Bắt buộc → must.", G7),
            mc("This ATM is ___ than the one at the hotel. (gần)", ["nearer", "more near", "nearest"], 0, "near → nearer than.", G6),
            fill("___ I pay by card here?", ["Can", "May"], "Xin phép: Can I…?", G7),
            fill("The rate here is ___ (good) than at the airport.", ["better"], "good → better than.", G6),
            translate("Bạn không được nói mật khẩu cho ai.", ["You mustn't tell anyone your password.", "You must not tell anyone your password.", "You mustn't tell your password to anyone.", "You must not tell your password to anyone."], "mustn't = cấm.", G7),
        ]),
        listening(k, "Đổi tiền ở ngân hàng", [
            ("A", "Good morning. Can I change dollars here?"), ("B", "Yes. How much do you want to change?"),
            ("A", "Two hundred dollars, please. What's the rate today?"), ("B", "Twenty-five thousand dong for one dollar."),
            ("A", "Is it better than at the airport?"), ("B", "Yes, banks are usually cheaper than the airport. Please sign here."),
        ], [
            numbers(k + "1", "Two hundred dollars, please.", "Anh ấy đổi bao nhiêu đô la?", "NUMBER", ["200"], "two hundred = 200."),
            numbers(k + "2", "Twenty-five thousand dong for one dollar.", "Một đô la đổi được bao nhiêu nghìn đồng?", "NUMBER", ["25", "25000", "25,000", "25.000"], "twenty-five thousand = 25.000."),
            true_false("The airport is cheaper than the bank.", False, "\"banks are usually cheaper than the airport\"", "listening"),
        ]),
        reading("Cách rút tiền ATM", "How to use this ATM. 1. Put your card in. 2. Choose your language. 3. Enter your PIN code. Never tell anyone your PIN. 4. Choose the amount. You can take out up to 5,000,000 VND a day. 5. Take your card and your money. The bank is open from 8 a.m. to 4:30 p.m., Monday to Friday.", [
            mc("Mỗi ngày rút được tối đa bao nhiêu?", ["5,000,000 VND", "500,000 VND", "50,000,000 VND"], 0, "\"up to 5,000,000 VND a day\"", "reading"),
            true_false("You can tell your PIN to the bank staff.", False, "\"Never tell anyone your PIN.\"", "reading"),
            mc("Ngân hàng mở cửa những ngày nào?", ["Monday to Friday", "every day", "Saturday and Sunday"], 0, "\"Monday to Friday\"", "reading"),
            true_false("You choose the language before you enter your PIN.", True, "Bước 2 chọn ngôn ngữ, bước 3 nhập PIN.", "reading"),
        ]),
        speaking(["Can I change dollars here?", "What's the rate today?", "Where do I sign?"]),
    ],
})

# Checkpoint lesson (week review) -------------------------------------------------
k = lesson("rv")
lessons.append({
    "slug": "a1-w3-review", "title": "Week 3 review",
    "description": "Ôn tập tuần 3: so sánh hơn / nhất, can / must / should, đi lại, khách sạn và sức khỏe.",
    "activities": [{
        "slug": "checkpoint", "type": "CHECKPOINT", "title": "Kiểm tra cuối tuần 3",
        "instructionsVi": "15 câu ngẫu nhiên từ các bài trong tuần, rồi một cuộc hội thoại với AI (không có gợi ý). "
                          "Cần 70 điểm và hoàn thành nhiệm vụ nói để mở tuần 4.",
        "scenarioSlug": "a1-w3-checkpoint",
        "content": {"fromActivityTypes": ["GRAMMAR", "LISTENING", "READING"], "itemCount": 15, "speaking": True},
        "exercises": [
            fill("The train is ___ (fast) than the bus.", ["faster"], "fast → faster than.", G6),
            fill("This is ___ (good) hotel in the city.", ["the best"], "good → the best.", G6),
            mc("You ___ smoke on the plane.", ["mustn't", "should", "can"], 0, "Cấm → mustn't.", G7),
            fill("You look tired. You ___ go to bed early.", ["should"], "Lời khuyên → should.", G7),
            fill("___ you tell me the way to the bank?", ["Can", "Could"], "Nhờ ai: Can / Could you…?", G7),
            error_correction("Hanoi is biger than Hue.", 2, ["bigger"], "big → bigger (gấp đôi g).", G6),
            mc("Which is ___? (expensive)", ["the most expensive", "the expensivest", "more expensive"], 0, "Tính từ dài → the most expensive.", G6),
            transform("You can park here.", "NEGATIVE", ["You can't park here.", "You cannot park here."], "can → can't.", G7),
            reorder("How much is a ticket to Hue?", "How much is + danh từ?", "vocab:travel"),
            translate("Rẽ phải ở góc đường.", ["Turn right at the corner."], "turn right = rẽ phải; at the corner = ở góc đường.", "vocab:directions"),
            dictation(k + "1", "Can I see your passport, please?", ""),
            numbers(k + "2", "The bus leaves at twenty past seven.", "Xe buýt chạy lúc mấy giờ?", "TIME", ["7:20"], "twenty past seven = 7:20."),
            minimal_pair(k + "1", [("think", "/θɪŋk/"), ("sink", "/sɪŋk/")], 0, "'think' có /θ/.", "sound:θ"),
            minimal_pair(k + "2", [("light", "/laɪt/"), ("right", "/raɪt/")], 1, "'right' có /r/.", "sound:r"),
        ],
    }],
})

# AI role-play (Tasks 67-70, D10): one scenario per situation lesson.
ROLE_PLAYS = {
    "a1-w3-l1-asking-directions": "a1-w3-asking-directions",
    "a1-w3-l2-taking-a-taxi": "a1-w3-taking-a-taxi",
    "a1-w3-l3-buying-tickets": "a1-w3-buying-tickets",
    "a1-w3-l4-airport-check-in": "a1-w3-airport-check-in",
    "a1-w3-l5-hotel-check-in": "a1-w3-hotel-check-in",
    "a1-w3-l6-train-bus-times": "a1-w3-train-bus-times",
    "a1-w3-l7-describing-symptoms": "a1-w3-describing-symptoms",
    "a1-w3-l8-buying-medicine": "a1-w3-buying-medicine",
    "a1-w3-l9-at-the-bank": "a1-w3-at-the-bank",
}

finish(lessons, {
    "slug": "a1-week-3", "title": "Week 3 · Travel and Health", "track": "CORE", "orderIndex": 3,
    "description": "Hỏi đường, taxi, mua vé, sân bay, khách sạn, giờ tàu xe, khám bệnh, mua thuốc và ngân hàng. Ngữ pháp: tính từ, so sánh hơn / nhất, can / must / should.",
    "icon": "plane", "emoji": "✈️", "color": "sky",
}, ROLE_PLAYS)
