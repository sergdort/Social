package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sergdort/Social/business/domain"
	"github.com/sergdort/Social/internal/store"
	"log"
	"math/rand"
	"time"
)

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers()

	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			log.Println("error creating user", err)
			_ = tx.Rollback()
			return
		}
	}
	_ = tx.Commit()

	posts := generatePosts(users)

	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("error creating post", err)
			return
		}
	}

	comments := generateComments(4, posts, users)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("error creating comment", err)
			return
		}
	}

	log.Println("seeding complete")
}

func generateComments(count int, posts []*store.Post, users []*domain.User) []*store.Comment {
	var comments []*store.Comment

	for _, post := range posts {
		postComments := getRandomComments(count) // Generate random comments for the post
		for _, postComment := range postComments {
			comments = append(comments, &store.Comment{
				PostID:  post.ID,
				UserID:  users[rand.Intn(len(users))].ID, // Pick a random user
				Content: postComment,
			})
		}
	}

	return comments
}

func generatePosts(users []*domain.User) []*store.Post {
	posts := make([]*store.Post, len(seedPosts))

	for i, seed := range seedPosts {
		posts[i] = &store.Post{
			UserID:  users[rand.Intn(len(users))].ID,
			Title:   seed.Title,
			Content: seed.Content,
			Tags:    getRandomTags(),
		}
	}

	return posts
}

func generateUsers() []*domain.User {
	users := make([]*domain.User, len(usernames))

	for i, username := range usernames {
		user := &domain.User{
			ID:        0,
			Username:  username,
			Email:     fmt.Sprintf("%s@example.com", username),
			CreatedAt: "",
			RoleID:    1,
		}
		_ = user.Password.Set("123123")
		users[i] = user
	}

	return users
}

var usernames = []string{
	"JonSnow",
	"KhaleesiStorm",
	"AryaTheFaceless",
	"Kingslayer47",
	"TyrionLionheart",
	"BranThreeEyed",
	"SansaWolfQueen",
	"LittlefingerMoves",
	"VarysTheSpider",
	"DavosSeaworth",
	"MelisandreFire",
	"TheHoundBurns",
	"NightKingChill",
	"SamwellTheWise",
	"RobertBaratheon",
	"DragonQueen23",
	"EuronCroweye",
	"GhostDirewolf",
	"GendryBaratheon",
	"OberynRedViper",
}

type SeedPost struct {
	Title   string
	Content string
}

var seedPosts = []SeedPost{
	{"Winter is Coming", "Winter is coming."},
	{"A Lannister Always Pays", "A Lannister always pays his debts."},
	{"The Night is Dark", "The night is dark and full of terrors."},
	{"Chaos is a Ladder", "Chaos isn’t a pit. Chaos is a ladder."},
	{"The Things I Do for Love", "The things I do for love."},
	{"Hold the Door", "Hold the door! Hold the door! Hodor!"},
	{"Valar Morghulis", "Valar Morghulis. All men must die."},
	{"Not Today", "What do we say to the God of Death? Not today."},
	{"The King in the North!", "The King in the North!"},
	{"Dracarys!", "Dracarys!"},
	{"Burn Them All", "Burn them all!"},
	{"I Drink and I Know Things", "That's what I do: I drink and I know things."},
	{"The North Remembers", "The North remembers."},
	{"A Lion Does Not Concern", "A lion does not concern itself with the opinion of sheep."},
	{"You Know Nothing", "You know nothing, Jon Snow."},
	{"The Pack Survives", "When the snows fall and the white winds blow, the lone wolf dies but the pack survives."},
	{"A Dragon is Not a Slave", "A dragon is not a slave."},
	{"Stick Them with the Pointy End", "Stick them with the pointy end."},
	{"Power is Power", "Power is power."},
	{"Fear Cuts Deeper", "Fear cuts deeper than swords."},
	{"Words are Wind", "Words are wind."},
	{"The Night's Watch Oath", "I am the sword in the darkness. I am the watcher on the walls."},
	{"There is No Middle Ground", "When you play the game of thrones, you win or you die. There is no middle ground."},
	{"A Man Must Die", "A man must die, but first he must be born."},
	{"Harder to Kill", "Some old wounds never truly heal, and bleed again at the slightest word."},
	{"Every Flight Begins", "Every flight begins with a fall."},
	{"Kill the Boy", "Kill the boy and let the man be born."},
	{"The Weak Crumble", "The strong have always taken from the weak."},
	{"A True King", "A true king is neither cruel nor fearful. He is strong and wise and brave."},
	{"A Mind Needs Books", "A mind needs books as a sword needs a whetstone."},
	{"If I Look Back", "If I look back, I am lost."},
	{"Never Forget", "The man who passes the sentence should swing the sword."},
	{"The Only Time", "The only time a man can be brave is when he is afraid."},
	{"The World is Built", "The world is built by killers."},
	{"The Iron Throne is Mine", "The Iron Throne is mine by rights."},
	{"A Girl Has No Name", "A girl has no name."},
	{"They Will Bend the Knee", "I will take what is mine. With fire and blood, I will take it."},
	{"The Mountain Will Fall", "Tell your father I'm here, and tell him the Lannisters aren't the only ones who pay their debts."},
	{"A Bastard's Duty", "You are a Stark. You may not have my name, but you have my blood."},
	{"Wherever Whores Go", "Wherever whores go…"},
	{"The Only Justice", "There is only one war that matters. The Great War. And it is here."},
	{"The King of Ashes", "I will not become a queen of ashes."},
	{"My Watch Begins", "Night gathers, and now my watch begins."},
	{"A Man is No One", "A man is no one, and that is something."},
	{"Oathkeeper", "The things we do for love."},
	{"A Direwolf's Loyalty", "The wolves will come again."},
	{"Honor Means Nothing", "Honor is a horse."},
	{"A Knife in the Dark", "The past is already written. The ink is dry."},
	{"War is Coming", "We are all human. Oh, we all do horrible things."},
	{"Victory Has a Price", "A lion does not fear sheep."},
	{"Old Gods and New", "There is only one god, and His name is Death."},
	{"The Blood of the First Men", "The First Men called us the children, but we were born long before them."},
	{"A Throne of Lies", "A ruler who hides behind paid executioners soon forgets what death is."},
	{"Happiness is a Lie", "The more you give a king, the more he wants."},
	{"A Lord's Burden", "We do not choose our destiny. We must do our duty, no? Great or small, we must do our duty."},
	{"Words are Cheap", "Some battles are won with swords and spears, others with quills and ravens."},
	{"Blood Demands Blood", "The blood of the First Men still flows in the veins of the Starks."},
	{"Iron Price or Gold Price?", "We do not sow."},
	{"All Will Kneel", "All kneel before the power of the king."},
	{"The Shield of Men", "The Night’s Watch is the shield that guards the realms of men."},
	{"Brotherhood of No Banners", "There is only one war, the war of life against death."},
	{"The Last of Her House", "Fire and blood."},
	{"The Dragon Wakes", "The Targaryens are closer to gods than men."},
	{"The Price of Loyalty", "I swore a vow to protect you."},
	{"No One is Safe", "The game of thrones is played by everyone."},
	{"The Lion’s Roar", "Hear me roar!"},
	{"A Crown for a King", "Gold crowns are cold things for the dead."},
	{"By Sword or Fire", "By sword or by fire, I will take back what is mine."},
	{"Lords of Winterfell", "A ruler who cannot protect his people is no ruler at all."},
	{"A Throne of Swords", "The Iron Throne is not a seat of comfort."},
	{"The Red Wedding", "The Lannisters send their regards."},
	{"Bastard No More", "The North will never forget."},
	{"The Lord of Light", "The night is dark and full of terrors."},
	{"A Cold Night Falls", "No one ever truly comes back from the dead."},
	{"The North is Strong", "The North remembers."},
}

var tags = []string{
	"GameOfThrones",
	"ASOIAF",
	"Winterfell",
	"Targaryen",
	"Stark",
	"Lannister",
	"NightWatch",
	"ValarMorghulis",
	"HouseOfTheDragon",
	"Dragons",
	"IronThrone",
	"KingsLanding",
	"WhiteWalkers",
	"RedWedding",
	"Dothraki",
	"FireAndBlood",
	"NorthRemembers",
	"FacelessMen",
	"BattleOfBastards",
	"DarkAndFullOfTerrors",
}

var comments = []string{
	// Happy comments
	"A fine tale! Even a Lannister can appreciate good words. – Tyrion Lannister",
	"You have my sword, my shield, and my thanks. – Jon Snow",
	"A feast for the mind as well as the belly! – Robert Baratheon",
	"Not bad for a knight in shining armor. – Brienne of Tarth",
	"Now this is a story worth singing about! – Jojen Reed",
	"A good story is like good wine—aged, deep, and worth savoring. – Olenna Tyrell",
	"Well said! The North remembers such wisdom. – Sansa Stark",
	"One does not need dragons to create fire with words. – Daenerys Targaryen",
	"You speak with the wisdom of the old gods. – Bran Stark",
	"I never thought I'd enjoy reading something so much. – Samwell Tarly",

	// Impressed comments
	"Now that is the work of a true master. – Tywin Lannister",
	"Few have words as sharp as Valyrian steel. – Oberyn Martell",
	"Your words hit like a warhammer, strong and true. – Gendry",
	"A tale worthy of the histories. – Maester Aemon",
	"Well played. You would do well in the game of thrones. – Littlefinger",
	"That was unexpected. You have my attention. – Varys",
	"Even the Red Keep’s scribes could not have done better. – Qyburn",
	"By the Seven, this is a tale worth telling! – Loras Tyrell",
	"If words were swords, you would be a knight. – Bronn",
	"You have the mind of a maester and the wit of a sellsword. – Davos Seaworth",

	// Angry comments
	"Burn them all! – Aerys II Targaryen",
	"You dare speak such words in my presence? – Cersei Lannister",
	"Take this down, take it down at once! – Joffrey Baratheon",
	"Enough! I have heard enough of your nonsense. – Stannis Baratheon",
	"I do not trust your words, nor do I trust you. – Roose Bolton",
	"You speak like a craven who has never seen battle. – Sandor Clegane",
	"This is an insult to my house! – Randyll Tarly",
	"Your words are as empty as your honor. – Euron Greyjoy",
	"Spare me your false courtesies. – Melisandre",
	"You think this is clever? I think not. – Walder Frey",

	// Sad comments
	"The world is full of broken things, and this just adds to it. – Jaime Lannister",
	"Some wounds never heal, and this has only made them worse. – Theon Greyjoy",
	"I had hoped for more, but hope is a fool’s game. – Catelyn Stark",
	"It is a cruel thing, and yet it does not surprise me. – Barristan Selmy",
	"The night is dark and full of disappointments. – Thoros of Myr",
	"Words are wind, and these words cut deep. – Arya Stark",
	"We are not meant for happiness, are we? – Lysa Arryn",
	"This is a bitter thing to read. – Edmure Tully",
	"It should not have ended this way. – Jorah Mormont",
	"Once, I believed in stories. Now, I know better. – Shireen Baratheon",
}

func getRandomTags() []string {
	rand.Seed(time.Now().UnixNano()) // Ensure randomness per run

	shuffled := make([]string, len(tags))
	copy(shuffled, tags)
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	return shuffled[:4] // Select first 4 shuffled elements
}

func getRandomComments(count int) []string {
	rand.Shuffle(len(comments), func(i, j int) { comments[i], comments[j] = comments[j], comments[i] })

	return comments[:count]
}
