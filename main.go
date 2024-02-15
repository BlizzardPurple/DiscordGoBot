package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // This line is changed
)

var (
	phrases             = []string{"fuck", "mf", "mofo", "shit", "randi", "mc", "bc", "bkl", "madarchod", "bhenchod", "hitler", "lode", "lawde", "chutiye", "chut", "lund", "bitch", "ass", "nigga", "fucked", "fucker", "fucking", "asshole", "nigger", "behenchod", "cunt", "twat", "bhadwe", "saale", "kamine", "chutiya"}
	userFizzes          = make(map[string]int)
	countries           = []string{"india", "us", "uk", "uae", "australia", "china", "brazil", "vietnam"}
	randSrc             = rand.NewSource(time.Now().UnixNano())
	randGen             = rand.New(randSrc)
	genderSpecificTerms = []string{"gentlemen", "ladies", "guys", "bros", "bois", "bros", "men", "women", "countrymen"}
)

type Answers struct {
	OriginChannelId string
	FavFood         string
	FavGame         string
	RecordId        int64
}

func (a *Answers) answersToMessageEmbed() discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Favorite Food",
			Value: a.FavFood,
		},
		{
			Name:  "Favorite Game",
			Value: a.FavGame,
		},
		{
			Name:  "Record Id",
			Value: strconv.FormatInt(a.RecordId, 10),
		},
	}

	return discordgo.MessageEmbed{
		Title:  "New responses!",
		Fields: fields,
		Color:  0x0000ff, // Blue color
	}
}

var responses map[string]Answers = map[string]Answers{}

func main() {
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // stopping database to close

	messageCreateHandler := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageCreate(s, m, db)
	}

	// Usinf the wrapped handler for discord server logic
	dg.AddHandler(messageCreateHandler)
	dg.AddHandler(onReady)
	dg.AddHandler(ReactionAddHandler)
	dg.AddHandler(ReactionRemoveHandler)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}
	defer dg.Close() //stopping discord session from closing

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// Ctrl-C toclose it
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}

/*
Function to create the final message is below
*/

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
	if m.Author.ID == s.State.User.ID {
		return // Message from self
	}

	//DM logic
	if m.GuildID == "" {
		handleDM(s, m, db)
	}

	//server logic
	fizzIncremented := false
	gendertalked := false
	admin := os.Getenv("ADMIN_ID")
	if m.Author.ID == admin {
		gendertalked = true
		fizzIncremented = true
	}

	content := strings.ToLower(m.Content)
	words := strings.Fields(content)

	//check for foul langugaes
	if !fizzIncremented {
		for _, phrase := range phrases {
			for _, word := range words {
				if word == phrase {
					//foul language found
					fizzIncremented = true
					userFizzes[m.Author.ID]++
					handleFizz(s, m)
					return
				}
			}
		}
	}

	// Check if the message contains any of the gender-specific terms.
	if !gendertalked {
		for _, term := range genderSpecificTerms {
			for _, word := range words {
				if word == term {
					//misgendered word found
					gendertalked = true
					handleMisgender(s, m, term)
					return
				}
			}
		}
	}

	if strings.HasPrefix(content, "hello everyone") || strings.HasPrefix(content, "!gobot hello") {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hello <@%s>!", m.Author.ID))
		if err != nil {
			log.Println("Error sending message: ", err)
		}
		return
	}

	if strings.HasPrefix(content, "!gobot country") {
		randomIndex := randGen.Intn(len(countries))
		s.ChannelMessageSend(m.ChannelID, countries[randomIndex])
		return
	}

	if strings.HasPrefix(content, "!gobot register") {
		UserPromptHandler(s, m)
		return
	}

	if strings.HasPrefix(content, "!gobot answers") {
		AnswersHandler(s, m, db)
		return
	}
}

func AnswersHandler(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
	spl := strings.Split(m.Content, " ")
	if len(spl) < 3 {
		s.ChannelMessageSend(m.ChannelID, "An ID must be provided. Ex: '!gobot answers  1'")
		return
	}

	id, err := strconv.Atoi(spl[2])
	if err != nil {
		log.Fatal(err)
	}

	var recordId int64
	var answerStr string
	var userId int64

	query := "SELECT * FROM discord_messages WHERE id = $1"
	row := db.QueryRow(query, id)
	err = row.Scan(&recordId, &answerStr, &userId)
	if err != nil {
		log.Fatal(err)
	}

	var answers Answers
	err = json.Unmarshal([]byte(answerStr), &answers)
	if err != nil {
		log.Fatal(err)
	}

	answers.RecordId = recordId
	embed := answers.answersToMessageEmbed()
	s.ChannelMessageSendEmbed(m.ChannelID, &embed)
}

func handleDM(s *discordgo.Session, m *discordgo.MessageCreate, db *sql.DB) {
	answers, ok := responses[m.ChannelID]
	if !ok {
		return
	}

	if answers.FavFood == "" {
		answers.FavFood = m.Content
		s.ChannelMessageSend(m.ChannelID, "Great! What's your favorite game now?")
		responses[m.ChannelID] = answers
		return
	} else {
		answers.FavGame = m.Content

		query := "INSERT INTO discord_messages (payload, user_id) VALUES ($1, $2) RETURNING id"
		jbytes, err := json.Marshal(answers)
		if err != nil {
			log.Fatal(err)
		}

		var lastInserted int64
		err = db.QueryRow(query, string(jbytes), m.ChannelID).Scan(&lastInserted)
		if err != nil {
			log.Fatal(err)
		}
		answers.RecordId = lastInserted

		s.ChannelMessageSend(m.ChannelID, "Thanks for the answer")
		// log.Printf("%s answers: %v, %v", m.ChannelID, answers.FavFood, answers.FavGame)
		embed := answers.answersToMessageEmbed()
		s.ChannelMessageSendEmbed(answers.OriginChannelId, &embed)

		delete(responses, m.ChannelID)
	}
}

func handleFizz(s *discordgo.Session, m *discordgo.MessageCreate) {
	badWordEmbed := &discordgo.MessageEmbed{
		Title:       "Mind your language",
		Description: fmt.Sprintf("<@%s> Foul Language warning = %d/10", m.Author.ID, userFizzes[m.Author.ID]),
		Color:       0xff0000,
	}

	if userFizzes[m.Author.ID] < 10 {
		s.ChannelMessageSendEmbed(m.ChannelID, badWordEmbed)
	}

	if userFizzes[m.Author.ID] >= 10 {
		err := s.GuildMemberDelete(m.GuildID, m.Author.ID)
		if err != nil {
			fmt.Println("Error kicking user:", err)
			return
		}

		userFizzes[m.Author.ID] = 0
		// fmt.Printf("Kicked user %s for reaching bad word count 10 or more.\n", m.Author.Username)

		kickNotificationEmbed := &discordgo.MessageEmbed{
			Title:       "Company policy violation",
			Description: fmt.Sprintf("Kicked user <@%s> for reaching foul language warning limit.", m.Author.ID),
			Color:       0xff0000,
		}

		s.ChannelMessageSendEmbed(m.ChannelID, kickNotificationEmbed)
	}
}

func handleMisgender(s *discordgo.Session, m *discordgo.MessageCreate, term string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Company policy violation",
		Description: fmt.Sprintf("Hey, <@%s>, \"%s\" is gender-specific. Kindly use a gender-neutral term like folks, everyone, ya'll, peeps, people, team, crew, pals, friends, beings, etc.", m.Author.ID, term),
		Color:       0xff0000,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func UserPromptHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// user channel created
	channel, err := s.UserChannelCreate((m.Author.ID))
	if err != nil {
		log.Panic(err)
	}

	if _, ok := responses[channel.ID]; !ok {
		responses[channel.ID] = Answers{
			OriginChannelId: m.ChannelID,
			FavFood:         "",
			FavGame:         "",
		}
		s.ChannelMessageSend(channel.ID, "Hey there! Here are some questions")
		s.ChannelMessageSend(channel.ID, "What is your favorite food?")
	} else {
		s.ChannelMessageSend(channel.ID, "We're still waiting... :)")
	}
}

var reactionMap sync.Map // Map to store the last reaction added by each user

func ReactionAddHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	lastReaction, ok := reactionMap.Load(r.UserID)
	if ok && lastReaction != r.Emoji.Name {
		return
	}

	reactionMap.Store(r.UserID, r.Emoji.Name)
	firenationid := os.Getenv("FIRE_NATION")
	waternationid := os.Getenv("WATER_NATION")
	if r.Emoji.Name == "🔥" {
		s.GuildMemberRoleAdd(r.GuildID, r.UserID, firenationid) // FirePeopleID
	} else if r.Emoji.Name == "💧" {
		s.GuildMemberRoleAdd(r.GuildID, r.UserID, waternationid) // WaterPeopleID
	}
}

func ReactionRemoveHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	firenationid := os.Getenv("FIRE_NATION")
	waternationid := os.Getenv("WATER_NATION")
	if r.Emoji.Name == "🔥" {
		s.GuildMemberRoleRemove(r.GuildID, r.UserID, firenationid) // FirePeopleID
	} else if r.Emoji.Name == "💧" {
		s.GuildMemberRoleRemove(r.GuildID, r.UserID, waternationid) // WaterPeopleID
	}
	reactionMap.Delete(r.UserID)
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	s.AddHandler(guildMemberAdd)
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Create an embed message to greet the new member.
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Company",
			IconURL: "", // Add your company's logo URL here if needed
		},
		Description: fmt.Sprintf("Welcome to the server, <@%s>! We're glad to have you here.", m.User.ID),
		Color:       0x00ff00, // Green color
	}

	// Send the embed message to the default channel of the server.
	_, err := s.ChannelMessageSendEmbed(m.GuildID, embed)
	if err != nil {
		fmt.Println("Error sending greeting message:", err)
	}
}
