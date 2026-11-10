"""Authoring source for the 9 A1 grammar theory cards (G1-G9).

Usage (from content/a1):  python3 build_grammar.py   →  a1-grammar.json
Import a1-grammar.json BEFORE the week bundles so GRAMMAR activities link to
their theory card. Exercises practise a card through skill tag "grammar:<slug>".
"""
import json


def ex(en, vi):
    return {"en": en, "vi": vi}


def mistake(wrong, right, note):
    return {"wrong": wrong, "right": right, "noteVi": note}


points = [
    {
        "slug": "g1-to-be",
        "title": "The verb to be: am / is / are",
        "titleVi": "Động từ to be",
        "summaryVi": "Dùng 'to be' để nói ai đó là ai, ở đâu, thế nào: tên, nghề nghiệp, quốc tịch, tuổi, cảm xúc. "
                     "Mỗi chủ ngữ đi với một dạng riêng. Khi nói, người ta thường rút gọn: I'm, you're, he's, she's, it's, we're, they're.",
        "formula": "I + am          → I am / I'm\n"
                   "He / She / It + is → She is / She's\n"
                   "You / We / They + are → They are / They're\n\n"
                   "Phủ định: S + am/is/are + not   (isn't, aren't)\n"
                   "Câu hỏi:  Am/Is/Are + S + …?   → Yes, I am. / No, she isn't.",
        "examples": [
            ex("I am a developer.", "Tôi là lập trình viên."),
            ex("She is from Da Nang.", "Cô ấy đến từ Đà Nẵng."),
            ex("We are colleagues.", "Chúng tôi là đồng nghiệp."),
            ex("He isn't busy today.", "Hôm nay anh ấy không bận."),
            ex("Are you married? — No, I'm not.", "Bạn đã kết hôn chưa? — Chưa."),
        ],
        "commonMistakes": [
            mistake("I is a teacher.", "I am a teacher.", "'I' luôn đi với 'am'."),
            mistake("He are my boss.", "He is my boss.", "He / She / It đi với 'is'."),
            mistake("I am work in a bank.", "I work in a bank.", "Không dùng 'am' trước động từ thường như 'work'."),
            mistake("She very busy.", "She is very busy.", "Tiếng Việt không cần 'là' trước tính từ, tiếng Anh thì bắt buộc có 'is'."),
            mistake("You are from Vietnam?", "Are you from Vietnam?", "Câu hỏi: đưa 'are' lên trước chủ ngữ."),
        ],
    },
    {
        "slug": "g2-nouns-articles",
        "title": "Nouns, plurals, a / an / the, this / these",
        "titleVi": "Danh từ, số nhiều, mạo từ và từ chỉ định",
        "summaryVi": "Danh từ đếm được có số ít và số nhiều. Danh từ số ít đếm được luôn cần a/an/the hoặc từ sở hữu (my, your…). "
                     "Dùng 'a' trước âm phụ âm, 'an' trước âm nguyên âm (an engineer, an hour). 'The' dùng cho thứ cả hai người đều biết. "
                     "This/these chỉ vật ở gần; that/those chỉ vật ở xa.",
        "formula": "Số nhiều: + s (books) · + es sau s, sh, ch, x (boxes, watches)\n"
                   "Bất quy tắc: child → children, man → men, woman → women, person → people\n"
                   "a + âm phụ âm (a teacher) · an + âm nguyên âm (an accountant)\n"
                   "this / that + số ít · these / those + số nhiều\n"
                   "Sở hữu: my, your, his, her, its, our, their + danh từ",
        "examples": [
            ex("I have two sisters and a brother.", "Tôi có hai chị gái và một anh trai."),
            ex("She is an engineer.", "Cô ấy là kỹ sư."),
            ex("These are my children.", "Đây là các con tôi."),
            ex("Is that your laptop?", "Kia có phải laptop của bạn không?"),
            ex("The meeting room is on the second floor.", "Phòng họp ở tầng hai."),
        ],
        "commonMistakes": [
            mistake("I am engineer.", "I am an engineer.", "Nghề nghiệp số ít cần a/an."),
            mistake("She is a accountant.", "She is an accountant.", "'accountant' bắt đầu bằng âm nguyên âm → an."),
            mistake("I have three child.", "I have three children.", "child → children (bất quy tắc)."),
            mistake("This books are new.", "These books are new.", "books số nhiều → these."),
            mistake("They are a doctors.", "They are doctors.", "Không dùng a/an với danh từ số nhiều."),
        ],
    },
    {
        "slug": "g3-present-simple",
        "title": "Present simple",
        "titleVi": "Thì hiện tại đơn",
        "summaryVi": "Dùng cho thói quen, lịch trình và sự thật: việc bạn làm hằng ngày, công việc, sở thích. "
                     "Với he/she/it, động từ thêm -s hoặc -es. Câu phủ định và câu hỏi dùng trợ động từ do/does, và động từ chính trở về nguyên mẫu.",
        "formula": "I / You / We / They + V        → I work in Hanoi.\n"
                   "He / She / It + V-s/es          → She works in Hanoi.\n"
                   "Phủ định: don't / doesn't + V    → He doesn't work on Sunday.\n"
                   "Câu hỏi:  Do / Does + S + V?     → Does she work here?\n"
                   "Trạng từ tần suất: always, usually, often, sometimes, never (đứng trước động từ thường)",
        "examples": [
            ex("I start work at nine.", "Tôi bắt đầu làm việc lúc chín giờ."),
            ex("She usually takes the bus.", "Cô ấy thường đi xe buýt."),
            ex("We don't work on Saturday.", "Chúng tôi không làm việc vào thứ Bảy."),
            ex("Does he drink coffee? — Yes, he does.", "Anh ấy có uống cà phê không? — Có."),
            ex("The shop opens at 8 a.m.", "Cửa hàng mở cửa lúc 8 giờ sáng."),
        ],
        "commonMistakes": [
            mistake("She work in a bank.", "She works in a bank.", "He/She/It → thêm -s vào động từ."),
            mistake("He doesn't works here.", "He doesn't work here.", "Sau doesn't, động từ trở về nguyên mẫu."),
            mistake("Do she like coffee?", "Does she like coffee?", "He/She/It dùng 'does'."),
            mistake("I am go to work by bus.", "I go to work by bus.", "Không dùng 'am' với động từ thường ở hiện tại đơn."),
            mistake("I go always to the gym.", "I always go to the gym.", "Trạng từ tần suất đứng trước động từ thường."),
        ],
    },
    {
        "slug": "g4-present-continuous",
        "title": "Present continuous",
        "titleVi": "Thì hiện tại tiếp diễn",
        "summaryVi": "Dùng cho việc đang xảy ra ngay lúc nói (now, at the moment) hoặc tạm thời trong giai đoạn này. "
                     "So sánh: I work in Hanoi (thường xuyên) với I'm working from home this week (tạm thời).",
        "formula": "S + am / is / are + V-ing\n"
                   "Phủ định: S + am/is/are + not + V-ing\n"
                   "Câu hỏi:  Am/Is/Are + S + V-ing?\n"
                   "Chính tả: make → making · sit → sitting · lie → lying",
        "examples": [
            ex("I'm having lunch now.", "Bây giờ tôi đang ăn trưa."),
            ex("She is talking to a customer.", "Cô ấy đang nói chuyện với khách hàng."),
            ex("They aren't working today.", "Hôm nay họ không làm việc."),
            ex("What are you doing?", "Bạn đang làm gì vậy?"),
        ],
        "commonMistakes": [
            mistake("I working now.", "I am working now.", "Không được bỏ am/is/are."),
            mistake("She is work at the moment.", "She is working at the moment.", "Sau is phải là V-ing."),
            mistake("I am liking coffee.", "I like coffee.", "Động từ chỉ cảm xúc, sở thích (like, love, know, want) thường không dùng -ing."),
            mistake("He is siting there.", "He is sitting there.", "sit → sitting (gấp đôi phụ âm cuối)."),
        ],
    },
    {
        "slug": "g5-there-is-prepositions",
        "title": "There is / there are + prepositions of place",
        "titleVi": "There is / There are và giới từ chỉ nơi chốn",
        "summaryVi": "Dùng 'there is / there are' để nói cái gì có ở đâu. 'There is' + danh từ số ít, 'there are' + danh từ số nhiều. "
                     "Kết hợp với giới từ: in, on, at, under, next to, between, near, behind, in front of.",
        "formula": "There is + a/an + danh từ số ít    → There is a bank near my house.\n"
                   "There are + danh từ số nhiều        → There are two meeting rooms.\n"
                   "Câu hỏi: Is there…? / Are there…?    · Phủ định: There isn't / There aren't\n"
                   "on + tầng (on the second floor) · in + phòng/thành phố · at + địa điểm cụ thể",
        "examples": [
            ex("There is a printer next to my desk.", "Có một máy in bên cạnh bàn của tôi."),
            ex("There are three people in my team.", "Nhóm tôi có ba người."),
            ex("Is there a toilet on this floor?", "Tầng này có nhà vệ sinh không?"),
            ex("The bank is between the cafe and the hotel.", "Ngân hàng ở giữa quán cà phê và khách sạn."),
        ],
        "commonMistakes": [
            mistake("It has a bank near here.", "There is a bank near here.", "Người Việt hay dịch 'có' thành 'have/has'; để nói vật tồn tại ở đâu, dùng 'there is'."),
            mistake("There is two rooms.", "There are two rooms.", "Danh từ số nhiều → there are."),
            mistake("My office is in the third floor.", "My office is on the third floor.", "Tầng: dùng 'on'."),
            mistake("The cafe is next the bank.", "The cafe is next to the bank.", "next to: không được bỏ 'to'."),
        ],
    },
    {
        "slug": "g6-adjectives-comparison",
        "title": "Adjectives, comparatives and superlatives",
        "titleVi": "Tính từ, so sánh hơn và so sánh nhất",
        "summaryVi": "Tính từ đứng trước danh từ (a big city) hoặc sau to be (The city is big). "
                     "So sánh hơn: tính từ ngắn + -er + than; tính từ dài: more + adj + than. "
                     "So sánh nhất: the + adj-est / the most + adj.",
        "formula": "Ngắn: cheap → cheaper than → the cheapest\n"
                   "Tận cùng -y: busy → busier → the busiest\n"
                   "Dài: expensive → more expensive than → the most expensive\n"
                   "Bất quy tắc: good → better → the best · bad → worse → the worst · far → farther → the farthest",
        "examples": [
            ex("Hanoi is a busy city.", "Hà Nội là một thành phố đông đúc."),
            ex("The bus is cheaper than a taxi.", "Xe buýt rẻ hơn taxi."),
            ex("This hotel is more expensive than that one.", "Khách sạn này đắt hơn khách sạn kia."),
            ex("It's the best restaurant in town.", "Đó là nhà hàng ngon nhất thị trấn."),
        ],
        "commonMistakes": [
            mistake("a city big", "a big city", "Tính từ đứng trước danh từ."),
            mistake("more cheap", "cheaper", "Tính từ ngắn thêm -er, không dùng more."),
            mistake("It is more better.", "It is better.", "better đã là so sánh hơn, không thêm more."),
            mistake("She is taller me.", "She is taller than me.", "So sánh hơn cần 'than'."),
        ],
    },
    {
        "slug": "g7-can-must-should",
        "title": "Modal verbs: can / must / should",
        "titleVi": "Động từ khuyết thiếu can / must / should",
        "summaryVi": "Can: khả năng hoặc xin phép / nhờ vả (Can I…? Can you…?). Must: bắt buộc (quy định). "
                     "Should: lời khuyên. Sau modal luôn là động từ nguyên mẫu, không thêm -s, không thêm 'to'.",
        "formula": "S + can / must / should + V (nguyên mẫu)\n"
                   "Phủ định: can't · mustn't (cấm) · shouldn't\n"
                   "Câu hỏi: Can I…? (xin phép) · Can you…? (nhờ) · Should I…? (hỏi ý kiến)\n"
                   "Lịch sự hơn: Could you…? · May I…?",
        "examples": [
            ex("Can you help me, please?", "Bạn giúp tôi được không?"),
            ex("You must wear a helmet.", "Bạn phải đội mũ bảo hiểm."),
            ex("You should drink more water.", "Bạn nên uống nhiều nước hơn."),
            ex("I can't come to the meeting.", "Tôi không thể đến cuộc họp."),
        ],
        "commonMistakes": [
            mistake("She cans swim.", "She can swim.", "Modal không thêm -s."),
            mistake("You should to rest.", "You should rest.", "Không có 'to' sau should / can / must."),
            mistake("He can speaks English.", "He can speak English.", "Sau can là động từ nguyên mẫu."),
            mistake("Do you can help me?", "Can you help me?", "Câu hỏi: đưa can lên đầu, không dùng do."),
        ],
    },
    {
        "slug": "g8-past-simple",
        "title": "Past simple",
        "titleVi": "Thì quá khứ đơn",
        "summaryVi": "Dùng cho việc đã xảy ra và kết thúc trong quá khứ (yesterday, last week, two days ago). "
                     "Động từ có quy tắc thêm -ed; động từ bất quy tắc phải học thuộc. Câu phủ định và câu hỏi dùng did, và động từ chính trở về nguyên mẫu. "
                     "Quá khứ của to be: was / were.",
        "formula": "S + V-ed / V2          → I worked late yesterday.\n"
                   "Phủ định: didn't + V     → She didn't come.\n"
                   "Câu hỏi:  Did + S + V?   → Did you call him?\n"
                   "to be: I/He/She/It was · You/We/They were\n"
                   "Bất quy tắc thường gặp: go→went · have→had · do→did · see→saw · get→got · make→made · take→took · come→came · buy→bought · eat→ate",
        "examples": [
            ex("I visited my parents last weekend.", "Cuối tuần trước tôi về thăm bố mẹ."),
            ex("We went to the beach.", "Chúng tôi đã đi biển."),
            ex("She didn't feel well yesterday.", "Hôm qua cô ấy thấy không khỏe."),
            ex("Did you finish the report? — Yes, I did.", "Bạn làm xong báo cáo chưa? — Xong rồi."),
            ex("The meeting was very long.", "Cuộc họp rất dài."),
        ],
        "commonMistakes": [
            mistake("I go to Da Lat last week.", "I went to Da Lat last week.", "Có 'last week' → động từ ở quá khứ."),
            mistake("She didn't went.", "She didn't go.", "Sau didn't, động từ trở về nguyên mẫu."),
            mistake("Did you saw him?", "Did you see him?", "Sau did, động từ trở về nguyên mẫu."),
            mistake("They was late.", "They were late.", "They → were."),
            mistake("I buyed a new phone.", "I bought a new phone.", "buy là động từ bất quy tắc: bought."),
        ],
    },
    {
        "slug": "g9-future",
        "title": "Future: will / be going to",
        "titleVi": "Tương lai: will và be going to",
        "summaryVi": "Be going to: kế hoạch đã định sẵn. Will: quyết định ngay lúc nói, lời hứa, đề nghị giúp đỡ, dự đoán. "
                     "Thì hiện tại tiếp diễn cũng dùng cho lịch hẹn đã sắp xếp (I'm meeting him at 3).",
        "formula": "S + am/is/are + going to + V   → I'm going to visit Hue next month.\n"
                   "S + will + V                    → I'll call you later.\n"
                   "Phủ định: won't + V · am/is/are not going to + V\n"
                   "Câu hỏi: Will you…? · Are you going to…?",
        "examples": [
            ex("I'm going to start a new course next week.", "Tuần sau tôi sẽ bắt đầu một khóa học mới."),
            ex("It's cold. I'll close the window.", "Trời lạnh. Để tôi đóng cửa sổ."),
            ex("I won't be late, I promise.", "Tôi sẽ không đến muộn, tôi hứa."),
            ex("What are you going to do this weekend?", "Cuối tuần này bạn định làm gì?"),
        ],
        "commonMistakes": [
            mistake("I will to call you.", "I will call you.", "Sau will là động từ nguyên mẫu, không có 'to'."),
            mistake("I going to travel.", "I'm going to travel.", "Không được bỏ am/is/are."),
            mistake("She wills come.", "She will come.", "will không thêm -s."),
            mistake("Tomorrow I go to Hue.", "Tomorrow I'm going to Hue.", "Kế hoạch tương lai: dùng be going to hoặc hiện tại tiếp diễn."),
        ],
    },
]

for index, point in enumerate(points, start=1):
    point["cefrLevel"] = "A1"
    point["orderIndex"] = index

bundle = {"version": 1, "courseSlug": "foundation-english", "cefrLevel": "A1", "grammarPoints": points, "topics": []}
with open("a1-grammar.json", "w", encoding="utf-8") as fh:
    json.dump(bundle, fh, ensure_ascii=False, indent=2)
print(f"grammar points={len(points)}")
