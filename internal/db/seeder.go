package db

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
)

func Seed(s store.Storage) error {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if _, err := s.Users.CreateWithInvitation(ctx, user, time.Hour*2); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	generatedPosts := generatePosts(users, len(posts))
	for _, post := range generatedPosts {
		if err := s.Posts.Create(ctx, post); err != nil {
			return fmt.Errorf("failed to create post: %w", err)
		}
	}

	generatedComments := generateComments(users, generatedPosts, len(generatedPosts)*3)
	for _, comment := range generatedComments {
		if err := s.Comments.Create(ctx, comment); err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
	}

	return nil
}

var names = []string{
	"Arlo", "Bryn", "Cael", "Dara", "Elio",
	"Fenn", "Gael", "Hale", "Idris", "Jove",
	"Kael", "Lior", "Mael", "Noel", "Orin",
	"Pax", "Quin", "Rael", "Soren", "Tael",
	"Ulan", "Vael", "Wren", "Xael", "Yael",
	"Zael", "Aiden", "Blaze", "Caden", "Drake",
	"Ember", "Flynn", "Greer", "Haven", "Ivor",
	"Jaden", "Kiran", "Lumen", "Maren", "Niall",
	"Oaken", "Piper", "Riven", "Sloan", "Taren",
	"Ulric", "Vance", "Wilder", "Xander", "Yoren",
	"Zoren", "Alder", "Birch", "Cedar", "Daven",
	"Eamon", "Faolan", "Gareth", "Hadyn", "Innes",
	"Jaryn", "Kelan", "Loryn", "Mavyn", "Naren",
	"Oswyn", "Pryn", "Rowan", "Stellan", "Taryn",
	"Ulwyn", "Voren", "Waryn", "Xoren", "Yoryn",
	"Zoren", "Aeron", "Boden", "Coren", "Doren",
	"Eoren", "Foren", "Goren", "Horen", "Ioren",
	"Joren", "Koren", "Loren", "Moren", "Noren",
	"Ooren", "Poren", "Roren", "Soren", "Toren",
	"Uoren", "Voren", "Woren", "Xoren", "Yoren",
}

var posts = []map[string]interface{}{
	{"title": "The Art of Minimalism", "content": "Minimalism is more than an aesthetic; it is a philosophy of living with intention. By reducing clutter, we create space for what truly matters.", "tags": []string{"lifestyle", "minimalism", "philosophy"}},
	{"title": "Mastering Go Concurrency", "content": "Goroutines and channels are the backbone of concurrent programming in Go. Understanding how to use them effectively can dramatically improve performance.", "tags": []string{"go", "programming", "concurrency"}},
	{"title": "The Science of Sleep", "content": "Sleep is not merely rest. It is an active process of repair, memory consolidation, and hormonal regulation that is essential for health.", "tags": []string{"health", "science", "sleep"}},
	{"title": "Urban Farming Revolution", "content": "Cities around the world are embracing rooftop gardens and vertical farms to bring food production closer to consumers and reduce carbon footprints.", "tags": []string{"farming", "urban", "sustainability"}},
	{"title": "Deep Dive into Docker", "content": "Containers have changed how we build and ship software. Docker provides a lightweight, portable environment that ensures consistency across development and production.", "tags": []string{"docker", "devops", "containers"}},
	{"title": "The Philosophy of Stoicism", "content": "Stoicism teaches that we cannot control external events, only our responses to them. This ancient wisdom remains remarkably relevant in modern life.", "tags": []string{"philosophy", "stoicism", "mindset"}},
	{"title": "Exploring Quantum Computing", "content": "Quantum computers leverage superposition and entanglement to solve problems that classical computers would take millions of years to crack.", "tags": []string{"quantum", "computing", "science"}},
	{"title": "The Joy of Sourdough", "content": "Baking sourdough is a meditative practice. The slow fermentation process develops complex flavours that commercial bread simply cannot replicate.", "tags": []string{"baking", "food", "sourdough"}},
	{"title": "Building REST APIs in Go", "content": "Go is an excellent language for building high performance REST APIs. Its standard library provides everything needed to get a server running in minutes.", "tags": []string{"go", "api", "backend"}},
	{"title": "The Psychology of Habits", "content": "Habits are formed through a loop of cue, routine, and reward. Understanding this cycle is the first step to building better behaviours.", "tags": []string{"psychology", "habits", "self-improvement"}},
	{"title": "Introduction to Machine Learning", "content": "Machine learning enables computers to learn from data without being explicitly programmed. It powers everything from recommendation engines to medical diagnostics.", "tags": []string{"ml", "ai", "technology"}},
	{"title": "The History of Jazz", "content": "Jazz emerged in New Orleans in the early 20th century, blending African rhythms with European harmonics to create one of Americas most original art forms.", "tags": []string{"music", "jazz", "history"}},
	{"title": "Mindful Meditation Practices", "content": "Meditation does not require hours of silence. Even five minutes of mindful breathing can reduce cortisol levels and improve focus throughout the day.", "tags": []string{"meditation", "mindfulness", "wellness"}},
	{"title": "Kubernetes for Beginners", "content": "Kubernetes orchestrates containerised applications at scale. While the learning curve is steep, its power to manage complex deployments is unmatched.", "tags": []string{"kubernetes", "devops", "containers"}},
	{"title": "The Economics of Attention", "content": "In the digital age, attention has become the scarcest resource. Companies compete fiercely for every second of your focus, reshaping how we think and behave.", "tags": []string{"economics", "technology", "society"}},
	{"title": "Photography Composition Tips", "content": "The rule of thirds, leading lines, and negative space are foundational principles that transform ordinary snapshots into compelling visual stories.", "tags": []string{"photography", "art", "composition"}},
	{"title": "Climate Change and Technology", "content": "From solar panels to carbon capture, technology is playing a crucial role in the fight against climate change, though policy remains equally important.", "tags": []string{"climate", "technology", "environment"}},
	{"title": "The Art of Public Speaking", "content": "Great public speakers are not born; they are made. Deliberate practice, storytelling, and controlled breathing are the foundations of confident delivery.", "tags": []string{"communication", "public-speaking", "skills"}},
	{"title": "PostgreSQL Performance Tuning", "content": "Indexes, query planning, and connection pooling are essential tools for squeezing maximum performance out of a PostgreSQL database under heavy load.", "tags": []string{"postgresql", "database", "performance"}},
	{"title": "The Rise of Remote Work", "content": "Remote work has shifted from a perk to an expectation for many knowledge workers. Companies that resist this shift risk losing top talent to more flexible competitors.", "tags": []string{"work", "remote", "culture"}},
	{"title": "Introduction to Rust", "content": "Rust offers memory safety without garbage collection, making it ideal for systems programming where performance and reliability are non-negotiable.", "tags": []string{"rust", "programming", "systems"}},
	{"title": "The Science of Nutrition", "content": "Macronutrients provide energy while micronutrients support cellular function. Understanding the balance between them is key to long term health.", "tags": []string{"nutrition", "health", "science"}},
	{"title": "Building a Personal Brand", "content": "Your personal brand is the story others tell about you when you are not in the room. Consistency, authenticity, and visibility are its three pillars.", "tags": []string{"branding", "career", "marketing"}},
	{"title": "The Future of Electric Vehicles", "content": "Battery technology improvements and expanding charging infrastructure are accelerating EV adoption, with internal combustion engines facing an uncertain future.", "tags": []string{"ev", "technology", "transport"}},
	{"title": "Design Patterns in Go", "content": "Design patterns like singleton, factory, and observer provide reusable solutions to common software engineering problems, even in idiomatic Go code.", "tags": []string{"go", "design-patterns", "programming"}},
	{"title": "The Neuroscience of Creativity", "content": "Creativity is not a gift but a skill rooted in the brains default mode network. Boredom and rest are not enemies of creativity; they fuel it.", "tags": []string{"neuroscience", "creativity", "psychology"}},
	{"title": "Exploring the Deep Ocean", "content": "Less than 20 percent of the ocean floor has been mapped. The deep sea remains one of Earths last great frontiers, teeming with undiscovered species.", "tags": []string{"ocean", "science", "exploration"}},
	{"title": "The Art of Negotiation", "content": "Effective negotiation is not about winning; it is about finding mutually beneficial outcomes. Listening more than you speak is the most underrated skill in any deal.", "tags": []string{"negotiation", "business", "skills"}},
	{"title": "GraphQL vs REST", "content": "GraphQL gives clients the power to request exactly the data they need, reducing over-fetching. REST remains simpler to cache and easier to reason about at scale.", "tags": []string{"graphql", "rest", "api"}},
	{"title": "The History of the Internet", "content": "From ARPANET to the World Wide Web, the internet evolved from a military communication tool to the backbone of modern civilisation in just a few decades.", "tags": []string{"internet", "history", "technology"}},
	{"title": "Strength Training Fundamentals", "content": "Progressive overload is the cornerstone of strength training. Consistently adding weight or reps over time is what drives muscle growth and increased performance.", "tags": []string{"fitness", "strength", "health"}},
	{"title": "The Ethics of Artificial Intelligence", "content": "As AI systems make increasingly consequential decisions, questions of bias, accountability, and transparency have moved from academic discussion to urgent policy debate.", "tags": []string{"ai", "ethics", "technology"}},
	{"title": "Exploring Nordic Cuisine", "content": "Nordic cuisine celebrates simplicity and local ingredients. Fermented fish, foraged herbs, and rye bread tell the story of a culture shaped by harsh winters.", "tags": []string{"food", "nordic", "culture"}},
	{"title": "Writing Clean Code", "content": "Clean code reads like well-written prose. Meaningful variable names, small functions, and clear intent reduce cognitive load and make maintenance a pleasure.", "tags": []string{"programming", "clean-code", "best-practices"}},
	{"title": "The Psychology of Money", "content": "Our relationship with money is shaped more by emotion and behaviour than by financial knowledge. Understanding your money story is the first step to changing it.", "tags": []string{"finance", "psychology", "money"}},
	{"title": "Space Tourism on the Horizon", "content": "Commercial space companies are making orbit accessible to private citizens for the first time. Within a decade, space tourism could become a genuine industry.", "tags": []string{"space", "tourism", "technology"}},
	{"title": "The Power of Daily Journaling", "content": "Writing daily forces clarity of thought. Whether processing emotions or planning goals, a journal is one of the most powerful and underused productivity tools.", "tags": []string{"journaling", "productivity", "wellness"}},
	{"title": "Redis Caching Strategies", "content": "Caching with Redis can reduce database load by orders of magnitude. Choosing between cache-aside, write-through, and write-behind depends on your consistency needs.", "tags": []string{"redis", "caching", "backend"}},
	{"title": "The Science of Happiness", "content": "Research consistently shows that relationships, purpose, and gratitude are stronger predictors of happiness than wealth or status. Hedonic adaptation humbles us all.", "tags": []string{"happiness", "psychology", "science"}},
	{"title": "Learning a Second Language", "content": "Immersion is the fastest path to fluency. Surrounding yourself with a language through media, conversation, and reading accelerates acquisition far beyond classroom study.", "tags": []string{"language", "learning", "education"}},
	{"title": "Microservices Architecture", "content": "Microservices decompose applications into small, independently deployable services. The trade-off is operational complexity in exchange for scalability and team autonomy.", "tags": []string{"microservices", "architecture", "backend"}},
	{"title": "The Art of Slow Travel", "content": "Slow travel prioritises depth over breadth. Spending weeks in one place rather than days in many allows genuine connection with local culture and people.", "tags": []string{"travel", "lifestyle", "culture"}},
	{"title": "Introduction to Blockchain", "content": "Blockchain is a distributed ledger that records transactions across many computers. Its decentralised nature makes tampering with records extremely difficult.", "tags": []string{"blockchain", "technology", "crypto"}},
	{"title": "The Benefits of Cold Exposure", "content": "Cold showers and ice baths activate the sympathetic nervous system, release norepinephrine, and have been shown to improve mood, alertness, and recovery.", "tags": []string{"health", "cold-exposure", "wellness"}},
	{"title": "Building Scalable Systems", "content": "Scalability begins with understanding your bottlenecks. Horizontal scaling, load balancing, and asynchronous processing are the primary levers available to engineers.", "tags": []string{"scalability", "systems", "engineering"}},
	{"title": "The Renaissance of Board Games", "content": "Modern board games have shed their reputation for simplicity. Designer games like Wingspan and Agricola offer deep strategic experiences for adults.", "tags": []string{"games", "culture", "leisure"}},
	{"title": "Functional Programming Concepts", "content": "Pure functions, immutability, and higher-order functions form the foundation of functional programming, leading to more predictable and testable code.", "tags": []string{"functional", "programming", "concepts"}},
	{"title": "The Gut Microbiome", "content": "The trillions of bacteria living in your gut influence mood, immunity, and metabolism. Diet is the single most powerful lever for shaping your microbiome.", "tags": []string{"health", "microbiome", "science"}},
	{"title": "Digital Minimalism", "content": "Intentionally reducing digital consumption frees attention for deeper work and meaningful relationships. Notifications are the enemy of sustained focus.", "tags": []string{"technology", "minimalism", "productivity"}},
	{"title": "The Art of Fermentation", "content": "Fermentation transforms simple ingredients through microbial activity. From kimchi to kombucha, fermented foods have sustained human civilisations for millennia.", "tags": []string{"food", "fermentation", "culture"}},
	{"title": "Event-Driven Architecture", "content": "Event-driven systems communicate through asynchronous messages, enabling loose coupling and high scalability. Kafka and RabbitMQ are the dominant brokers in this space.", "tags": []string{"architecture", "events", "backend"}},
	{"title": "The History of Philosophy", "content": "From Socrates to Wittgenstein, philosophy has grappled with questions of existence, knowledge, and ethics. Its history is a conversation across centuries.", "tags": []string{"philosophy", "history", "education"}},
	{"title": "Cycling as a Lifestyle", "content": "Cycling combines exercise, transport, and adventure. Cities investing in cycling infrastructure see measurable improvements in public health and reduced congestion.", "tags": []string{"cycling", "fitness", "lifestyle"}},
	{"title": "WebSockets in Practice", "content": "WebSockets enable full-duplex communication between client and server, making them ideal for real-time applications like chat, live dashboards, and multiplayer games.", "tags": []string{"websockets", "backend", "real-time"}},
	{"title": "The Power of Compounding", "content": "Compounding works in finance, skills, and relationships. Small consistent improvements accumulate into transformative results over years and decades.", "tags": []string{"finance", "growth", "mindset"}},
	{"title": "Wilderness Survival Skills", "content": "Fire, shelter, water, and food — in that order. Wilderness survival is about prioritising immediate threats and remaining calm when plans fall apart.", "tags": []string{"survival", "outdoors", "skills"}},
	{"title": "The Science of Exercise", "content": "Exercise triggers the release of BDNF, a protein that promotes neuronal growth. Regular physical activity is one of the most effective interventions for mental health.", "tags": []string{"exercise", "science", "health"}},
	{"title": "Observability in Production", "content": "Logs, metrics, and traces form the three pillars of observability. Without them, debugging production issues is like navigating blindfolded through a maze.", "tags": []string{"devops", "observability", "engineering"}},
	{"title": "The Philosophy of Time", "content": "Is time a fundamental feature of reality or an emergent property of change? Physicists and philosophers continue to debate the nature of time without resolution.", "tags": []string{"philosophy", "time", "science"}},
	{"title": "Home Coffee Brewing Guide", "content": "Grind size, water temperature, and extraction time are the three variables that determine coffee quality. A burr grinder is the single best investment for home brewing.", "tags": []string{"coffee", "food", "guide"}},
	{"title": "The Future of Work", "content": "Automation will displace routine jobs while creating demand for creative, interpersonal, and technical skills. Lifelong learning is no longer optional.", "tags": []string{"work", "future", "technology"}},
	{"title": "Introduction to TypeScript", "content": "TypeScript adds static typing to JavaScript, catching errors at compile time rather than runtime. For large codebases, the productivity gains are substantial.", "tags": []string{"typescript", "javascript", "programming"}},
	{"title": "The Art of Letter Writing", "content": "A handwritten letter carries weight that email cannot. The deliberate effort of putting pen to paper communicates care in a way digital messages rarely achieve.", "tags": []string{"writing", "communication", "culture"}},
	{"title": "Renewable Energy Trends", "content": "Solar and wind costs have fallen faster than any model predicted. Renewables are now the cheapest source of new electricity generation in most of the world.", "tags": []string{"energy", "sustainability", "technology"}},
	{"title": "Database Indexing Explained", "content": "An index trades storage space for query speed. Understanding B-tree and hash indexes, and knowing when not to index, is fundamental to database performance.", "tags": []string{"database", "indexing", "backend"}},
	{"title": "The Benefits of Reading Fiction", "content": "Fiction builds empathy by immersing readers in perspectives unlike their own. Research shows regular fiction readers score higher on tests of social cognition.", "tags": []string{"reading", "fiction", "education"}},
	{"title": "Rock Climbing for Beginners", "content": "Rock climbing builds grip strength, problem-solving, and mental resilience. Indoor climbing gyms offer a safe and social environment to learn the basics.", "tags": []string{"climbing", "fitness", "outdoors"}},
	{"title": "API Rate Limiting Strategies", "content": "Token bucket and sliding window algorithms are the most common approaches to rate limiting. Choosing the right strategy depends on your traffic patterns and fairness requirements.", "tags": []string{"api", "rate-limiting", "backend"}},
	{"title": "The History of Typography", "content": "From Gutenbergs movable type to variable fonts, typography has shaped how humans read and communicate for over five centuries.", "tags": []string{"typography", "design", "history"}},
	{"title": "Plant Based Diet Guide", "content": "A well-planned plant-based diet provides all essential nutrients except B12. The environmental benefits of reducing animal product consumption are substantial.", "tags": []string{"diet", "health", "environment"}},
	{"title": "Building CLI Tools in Go", "content": "Go produces single binary executables that make distributing CLI tools trivial. The cobra library provides a robust framework for building complex command-line interfaces.", "tags": []string{"go", "cli", "programming"}},
	{"title": "The Art of Saying No", "content": "Every yes is implicitly a no to something else. Learning to decline requests that do not align with your priorities is a skill that compounds over a lifetime.", "tags": []string{"productivity", "mindset", "self-improvement"}},
	{"title": "Astronomy for Beginners", "content": "A pair of binoculars is enough to observe the moons of Jupiter and the craters of our own Moon. The night sky rewards those patient enough to look up.", "tags": []string{"astronomy", "science", "outdoors"}},
	{"title": "OAuth 2.0 Explained", "content": "OAuth 2.0 is an authorisation framework that allows third-party applications to access user data without exposing credentials. Understanding its flows is essential for any web developer.", "tags": []string{"oauth", "security", "backend"}},
	{"title": "The Craft of Screenwriting", "content": "A screenplay is a blueprint for a film, not the film itself. Economy of language, visual storytelling, and subtext distinguish great scripts from merely competent ones.", "tags": []string{"writing", "film", "craft"}},
	{"title": "Intermittent Fasting Science", "content": "Intermittent fasting works primarily through caloric restriction, though metabolic switching to ketones during fasting periods may offer additional cognitive benefits.", "tags": []string{"health", "fasting", "science"}},
	{"title": "The Power of Open Source", "content": "Open source software has become the foundation of the modern internet. The collaborative model of development consistently produces more secure and innovative software.", "tags": []string{"open-source", "technology", "community"}},
	{"title": "Landscape Photography Tips", "content": "Golden hour light transforms ordinary landscapes into extraordinary images. Patience, composition, and understanding weather patterns separate good landscape photographers from great ones.", "tags": []string{"photography", "landscape", "art"}},
	{"title": "Message Queue Patterns", "content": "Message queues decouple producers from consumers, enabling asynchronous processing and improving resilience. Dead letter queues are essential for handling failed messages gracefully.", "tags": []string{"messaging", "queues", "architecture"}},
	{"title": "The Science of Longevity", "content": "Caloric restriction, exercise, sleep, and social connection are the four lifestyle factors most consistently associated with longer healthspan in the scientific literature.", "tags": []string{"longevity", "health", "science"}},
	{"title": "Learning Chess as an Adult", "content": "Chess rewards pattern recognition and strategic thinking. Adult learners who study endgames before openings progress faster than those who do the opposite.", "tags": []string{"chess", "games", "learning"}},
	{"title": "Zero Downtime Deployments", "content": "Blue-green deployments, canary releases, and feature flags are the primary strategies for deploying software without interrupting users. Each involves different trade-offs in complexity and cost.", "tags": []string{"devops", "deployment", "engineering"}},
	{"title": "The History of Architecture", "content": "Architecture reflects the values, technology, and resources of its time. From Gothic cathedrals to modernist glass towers, buildings tell the story of civilisation.", "tags": []string{"architecture", "history", "art"}},
	{"title": "Sourdough Starter Guide", "content": "A sourdough starter is a living culture of wild yeast and bacteria. Feeding it regularly with flour and water maintains the microbial balance that leavens bread.", "tags": []string{"baking", "sourdough", "food"}},
	{"title": "The Art of Doing Nothing", "content": "The Italian concept of dolce far niente celebrates the sweetness of doing nothing. Rest and idleness are not laziness; they are necessary conditions for creativity.", "tags": []string{"lifestyle", "wellness", "culture"}},
	{"title": "Go Error Handling Patterns", "content": "Go treats errors as values, not exceptions. Wrapping errors with context using fmt.Errorf and errors.Is provides a clear chain of responsibility without panic.", "tags": []string{"go", "errors", "programming"}},
	{"title": "The Science of Addiction", "content": "Addiction hijacks the brains reward circuitry, creating compulsive behaviour despite harmful consequences. Understanding the neuroscience of addiction challenges moral judgements about it.", "tags": []string{"science", "addiction", "psychology"}},
	{"title": "Urban Cycling Infrastructure", "content": "Protected bike lanes reduce cycling fatalities dramatically. Cities that invest in cycling infrastructure see increased ridership, improved public health, and reduced traffic congestion.", "tags": []string{"cycling", "urban", "infrastructure"}},
	{"title": "Distributed Tracing with OpenTelemetry", "content": "OpenTelemetry provides a vendor-neutral standard for collecting traces, metrics, and logs. Instrumenting your services with it from the start saves significant pain later.", "tags": []string{"observability", "tracing", "devops"}},
	{"title": "The Philosophy of Language", "content": "Does language shape thought, or does thought shape language? The Sapir-Whorf hypothesis suggests our native tongue influences how we perceive reality itself.", "tags": []string{"philosophy", "language", "psychology"}},
	{"title": "Fermented Drinks Around the World", "content": "From Ethiopian tej to Japanese sake, every culture has developed fermented beverages from local ingredients. Each reflects the microbial terroir of its region.", "tags": []string{"drinks", "fermentation", "culture"}},
	{"title": "Go Generics in Practice", "content": "Generics arrived in Go 1.18, enabling type-safe data structures and algorithms without code duplication. Use them where they reduce boilerplate, not to over-engineer simple solutions.", "tags": []string{"go", "generics", "programming"}},
	{"title": "The Art of Listening", "content": "Most people listen to respond, not to understand. True listening requires suspending your own narrative and being genuinely curious about the other persons experience.", "tags": []string{"communication", "skills", "relationships"}},
	{"title": "Marathon Training Basics", "content": "Building an aerobic base over months, not weeks, is the key to marathon success. Most runners run their long runs too fast and their easy runs too hard.", "tags": []string{"running", "marathon", "fitness"}},
	{"title": "Service Mesh Architecture", "content": "A service mesh like Istio or Linkerd handles cross-cutting concerns like mTLS, retries, and circuit breaking at the infrastructure layer, freeing developers from implementing them in code.", "tags": []string{"service-mesh", "architecture", "devops"}},
	{"title": "The History of Money", "content": "Money evolved from barter to commodity money to fiat currency. Each transition reflected changing social trust and the need for more portable, divisible stores of value.", "tags": []string{"money", "history", "economics"}},
	{"title": "Indoor Plant Care Guide", "content": "Light and watering frequency are the two variables that determine plant health. Most indoor plants die from overwatering, not underwatering, so err on the side of dry.", "tags": []string{"plants", "home", "guide"}},
	{"title": "The Science of Memory", "content": "Memory is not a recording but a reconstruction. Each time we recall a memory we alter it slightly, making eyewitness testimony far less reliable than courts historically assumed.", "tags": []string{"memory", "neuroscience", "psychology"}},
	{"title": "CI/CD Pipeline Best Practices", "content": "A fast CI pipeline is a competitive advantage. Parallelising tests, caching dependencies, and running only affected tests are the primary levers for reducing build times.", "tags": []string{"ci-cd", "devops", "engineering"}},
	{"title": "The Art of Improvisation", "content": "Yes, and — the foundational rule of improv comedy — is also a powerful principle for collaboration. Building on others ideas rather than blocking them unlocks collective creativity.", "tags": []string{"improv", "creativity", "communication"}},
	{"title": "Understanding Macroeconomics", "content": "GDP, inflation, and unemployment form the core metrics of macroeconomic health. Central banks use interest rates as their primary tool for influencing all three.", "tags": []string{"economics", "macro", "finance"}},
	{"title": "Trail Running Guide", "content": "Trail running builds proprioception and mental toughness that road running cannot. Start on gentle trails and prioritise time on feet over pace.", "tags": []string{"running", "trails", "fitness"}},
	{"title": "HTTP/2 and HTTP/3 Explained", "content": "HTTP/2 introduced multiplexing to eliminate head-of-line blocking. HTTP/3 goes further by replacing TCP with QUIC, dramatically improving performance on lossy networks.", "tags": []string{"http", "networking", "backend"}},
	{"title": "The Philosophy of Mind", "content": "The hard problem of consciousness asks why physical processes in the brain give rise to subjective experience. No scientific theory has yet come close to resolving it.", "tags": []string{"philosophy", "consciousness", "mind"}},
	{"title": "Foraging Wild Edibles", "content": "Wild garlic, nettles, and elderflower are among the easiest plants to identify and harvest safely. Always use multiple identification features and never eat anything you are not certain of.", "tags": []string{"foraging", "food", "outdoors"}},
	{"title": "Go Testing Best Practices", "content": "Table-driven tests keep test cases organised and easy to extend. The testify library adds expressive assertions, while httptest makes handler testing clean and fast.", "tags": []string{"go", "testing", "programming"}},
	{"title": "The Science of Flow States", "content": "Flow occurs when challenge and skill are in balance. Csikszentmihalyi identified autotelic experience as the key to intrinsic motivation and sustained engagement.", "tags": []string{"psychology", "flow", "productivity"}},
	{"title": "Building a Home Studio", "content": "Acoustic treatment matters more than expensive microphones. Absorbing early reflections with panels and bass traps transforms any room into a usable recording space.", "tags": []string{"music", "studio", "audio"}},
	{"title": "The Future of Databases", "content": "NewSQL databases attempt to combine the ACID guarantees of relational databases with the horizontal scalability of NoSQL systems. The trade-offs remain complex.", "tags": []string{"database", "future", "technology"}},
	{"title": "Wild Swimming", "content": "Swimming in rivers, lakes, and seas connects you to the natural world in a way pools cannot replicate. Cold water immersion has well-documented mental health benefits.", "tags": []string{"swimming", "outdoors", "wellness"}},
	{"title": "The Art of Delegation", "content": "Delegation is not abdication. Effective delegation involves clear expectations, appropriate authority, and feedback loops — not simply assigning tasks and disappearing.", "tags": []string{"leadership", "management", "skills"}},
	{"title": "Biomechanics of Running", "content": "Cadence, foot strike, and hip extension are the three biomechanical variables most correlated with running efficiency and injury prevention.", "tags": []string{"running", "biomechanics", "science"}},
	{"title": "gRPC vs REST", "content": "gRPC uses Protocol Buffers and HTTP/2 to deliver significantly faster inter-service communication than REST. It excels in internal microservice communication where browser compatibility is not required.", "tags": []string{"grpc", "rest", "backend"}},
	{"title": "The History of Writing", "content": "Writing was independently invented at least three times in human history. Sumerian cuneiform, Egyptian hieroglyphs, and Mesoamerican scripts each arose from the need to record trade.", "tags": []string{"writing", "history", "culture"}},
	{"title": "Beekeeping for Beginners", "content": "A single hive can house 60,000 bees and produce 27 kilograms of honey per year. Beekeeping teaches patience, observation, and respect for complex natural systems.", "tags": []string{"beekeeping", "nature", "hobby"}},
	{"title": "Connection Pooling in Go", "content": "database/sql manages a pool of connections automatically, but configuring MaxOpenConns and MaxIdleConns correctly for your workload prevents both exhaustion and unnecessary overhead.", "tags": []string{"go", "database", "performance"}},
	{"title": "The Psychology of Procrastination", "content": "Procrastination is primarily an emotional regulation problem, not a time management one. We avoid tasks associated with negative emotions, not tasks that are simply difficult.", "tags": []string{"psychology", "procrastination", "productivity"}},
	{"title": "Bread Baking Fundamentals", "content": "Gluten development through mixing and fermentation gives bread its structure. Understanding the role of hydration, salt, and time unlocks endless variation from simple ingredients.", "tags": []string{"baking", "bread", "food"}},
	{"title": "Infrastructure as Code", "content": "Terraform and Pulumi allow infrastructure to be versioned, reviewed, and reproduced. Manual infrastructure is the enemy of consistency, auditability, and disaster recovery.", "tags": []string{"infrastructure", "devops", "terraform"}},
	{"title": "The Science of Colour", "content": "Colour is not a property of objects but a perceptual experience created by the interaction of light, surface, and visual system. Two people never see exactly the same colour.", "tags": []string{"science", "colour", "perception"}},
	{"title": "Sea Kayaking Adventures", "content": "Sea kayaking combines physical challenge with access to coastlines unreachable by foot or boat. Reading tides and weather is as important as paddling technique.", "tags": []string{"kayaking", "outdoors", "adventure"}},
	{"title": "The Art of Feedback", "content": "Effective feedback is specific, timely, and focused on behaviour rather than identity. The sandwich method is well-intentioned but usually dilutes the message.", "tags": []string{"feedback", "communication", "leadership"}},
	{"title": "Understanding Tidal Forces", "content": "Tides arise from the gravitational differential between the near and far sides of Earth relative to the Moon. The Sun contributes about a third of total tidal force.", "tags": []string{"science", "tides", "astronomy"}},
	{"title": "Middleware Patterns in Go", "content": "HTTP middleware chains in Go use closure and the http.Handler interface to compose request processing logic cleanly. Logging, authentication, and rate limiting all belong at this layer.", "tags": []string{"go", "middleware", "backend"}},
	{"title": "The Philosophy of Ethics", "content": "Consequentialism judges actions by outcomes. Deontology judges by adherence to rules. Virtue ethics asks what a person of good character would do. Each captures something the others miss.", "tags": []string{"ethics", "philosophy", "morality"}},
	{"title": "Mushroom Cultivation at Home", "content": "Oyster mushrooms can be grown on straw or coffee grounds with minimal equipment. The mycelium colonises the substrate within weeks before fruiting bodies emerge.", "tags": []string{"mushrooms", "food", "hobby"}},
	{"title": "Chaos Engineering", "content": "Chaos engineering deliberately introduces failures into production systems to discover weaknesses before they cause outages. Netflix pioneered the practice with Chaos Monkey.", "tags": []string{"chaos-engineering", "devops", "reliability"}},
	{"title": "The History of Mathematics", "content": "Mathematics has been discovered and rediscovered across cultures. Babylonian algebra, Greek geometry, and Indian numerals each contributed to the global edifice of mathematical knowledge.", "tags": []string{"mathematics", "history", "science"}},
	{"title": "Cold Brew Coffee at Home", "content": "Cold brew is made by steeping coarsely ground coffee in cold water for 12 to 24 hours. The result is a concentrate that is smoother and less acidic than hot-brewed coffee.", "tags": []string{"coffee", "food", "guide"}},
	{"title": "Building Resilient Systems", "content": "Circuit breakers, retries with exponential backoff, and bulkheads are the primary patterns for building systems that degrade gracefully under failure conditions.", "tags": []string{"resilience", "architecture", "engineering"}},
	{"title": "The Art of Botanical Illustration", "content": "Botanical illustration combines scientific accuracy with artistic beauty. Every vein, texture, and seed must be rendered faithfully to serve its purpose as a scientific record.", "tags": []string{"art", "botany", "illustration"}},
	{"title": "Strength Training for Runners", "content": "Single-leg exercises like Bulgarian split squats build the stability and power that reduce injury risk and improve running economy at all distances.", "tags": []string{"strength", "running", "fitness"}},
	{"title": "DNS Explained", "content": "DNS translates human-readable domain names into IP addresses. Understanding TTLs, record types, and resolver caching is essential for diagnosing web connectivity issues.", "tags": []string{"dns", "networking", "backend"}},
	{"title": "The Psychology of Leadership", "content": "Transformational leaders inspire through vision and values rather than control. Psychological safety — the belief that speaking up carries no punishment — is the foundation of high-performing teams.", "tags": []string{"leadership", "psychology", "management"}},
	{"title": "Campfire Cooking Techniques", "content": "Cast iron, foil packets, and skewers are the three pillars of campfire cooking. Managing coal temperature rather than flame height is the key to consistent results outdoors.", "tags": []string{"cooking", "outdoors", "food"}},
	{"title": "Load Testing Your API", "content": "k6 and Locust are excellent tools for simulating realistic traffic patterns. Load testing reveals throughput limits and latency degradation before real users do.", "tags": []string{"testing", "api", "performance"}},
	{"title": "The History of Film", "content": "Cinema began as a fairground novelty and became the defining art form of the 20th century. The transition from silent to sound film eliminated many careers overnight.", "tags": []string{"film", "history", "culture"}},
	{"title": "Yoga for Athletes", "content": "Yoga improves the mobility, breathing efficiency, and body awareness that enhance performance in almost every sport. Hip openers and thoracic rotations are particularly valuable for runners.", "tags": []string{"yoga", "fitness", "wellness"}},
	{"title": "Go Modules Explained", "content": "Go modules replaced GOPATH-based dependency management with a versioned, reproducible system. Understanding go.sum and module proxies prevents supply chain surprises.", "tags": []string{"go", "modules", "programming"}},
	{"title": "The Science of Climate Systems", "content": "The Atlantic Meridional Overturning Circulation distributes heat around the planet. Its weakening due to freshwater influx from melting ice represents one of the most concerning climate tipping points.", "tags": []string{"climate", "science", "environment"}},
	{"title": "Bookbinding as a Craft", "content": "Hand-bound books have a quality and intimacy that printed volumes cannot match. Coptic stitch binding leaves the spine open, allowing books to lie flat when open.", "tags": []string{"craft", "bookbinding", "art"}},
	{"title": "Secrets of Great UX Writing", "content": "UX writing is the craft of guiding users through interfaces with clear, concise, and human language. Every button label and error message is an opportunity to build trust.", "tags": []string{"ux", "writing", "design"}},
	{"title": "The Art of Pickling", "content": "Pickling preserves vegetables through acidity, either from vinegar or lactic acid fermentation. The brine ratio determines both shelf life and flavour intensity.", "tags": []string{"food", "pickling", "preservation"}},
	{"title": "TLS and HTTPS Explained", "content": "TLS encrypts data in transit using a combination of asymmetric key exchange and symmetric encryption. Certificate transparency logs have made it significantly harder to issue fraudulent certificates.", "tags": []string{"security", "tls", "networking"}},
	{"title": "The Joy of Birdwatching", "content": "Birdwatching trains attention and patience. A field guide, binoculars, and a local patch are all you need to discover that the natural world is richer than most people realise.", "tags": []string{"birdwatching", "nature", "outdoors"}},
	{"title": "Dependency Injection in Go", "content": "Passing dependencies explicitly rather than using globals makes code testable and modular. Wire and Dig automate the wiring for large applications without sacrificing clarity.", "tags": []string{"go", "dependency-injection", "programming"}},
	{"title": "The Philosophy of Free Will", "content": "Compatibilists argue that free will and determinism can coexist. Hard determinists say our choices are the inevitable result of prior causes. The debate remains unresolved.", "tags": []string{"philosophy", "free-will", "consciousness"}},
	{"title": "Ceramics for Beginners", "content": "Centering clay on the wheel is the most challenging skill in pottery and the foundation of everything else. Patience and consistent pressure are more important than strength.", "tags": []string{"ceramics", "craft", "art"}},
	{"title": "Monitoring with Prometheus", "content": "Prometheus scrapes metrics from instrumented services and stores them as time series. Its query language PromQL enables powerful alerting and dashboarding with Grafana.", "tags": []string{"prometheus", "monitoring", "devops"}},
	{"title": "The History of Tea", "content": "Tea drinking originated in China millennia ago before spreading along trade routes to transform cultures worldwide. The British tea ritual reshaped agricultural systems across Asia.", "tags": []string{"tea", "history", "culture"}},
	{"title": "Building Empathy Through Travel", "content": "Exposure to different ways of living dismantles assumptions about how the world must work. Travel is one of the most reliable methods for expanding moral imagination.", "tags": []string{"travel", "empathy", "culture"}},
	{"title": "Concurrency Patterns in Go", "content": "Fan-out, fan-in, pipelines, and worker pools are the fundamental concurrency patterns in Go. Each solves a different class of parallelism problem with goroutines and channels.", "tags": []string{"go", "concurrency", "patterns"}},
	{"title": "The Science of Soil", "content": "A teaspoon of healthy soil contains more microorganisms than there are people on Earth. Soil health determines agricultural productivity and is a critical carbon sink.", "tags": []string{"soil", "science", "environment"}},
	{"title": "The Art of Watercolour", "content": "Watercolour rewards transparency and spontaneity. Working light to dark and leaving white space are the fundamental disciplines that separate skilled watercolourists from beginners.", "tags": []string{"art", "watercolour", "painting"}},
	{"title": "Secrets of Effective Retrospectives", "content": "The best retrospectives produce one or two concrete actions, not long lists of complaints. Safety to speak honestly is more important than any specific format or facilitation technique.", "tags": []string{"agile", "retrospectives", "teams"}},
	{"title": "The History of Cartography", "content": "Maps are never neutral. Every cartographic choice — projection, scale, what to include — reflects the assumptions and priorities of its maker and era.", "tags": []string{"maps", "history", "geography"}},
	{"title": "Strength in Stillness", "content": "Stillness is not passivity. The ability to sit with discomfort, resist distraction, and act from a place of calm is among the most powerful capacities a person can develop.", "tags": []string{"mindfulness", "philosophy", "wellbeing"}},
}

func generateUsers(n int) []*model.User {
	users := make([]*model.User, n)
	for i := 0; i < n; i++ {
		users[i] = &model.User{
			FirstName: names[rand.Intn(len(names))],
			LastName:  names[rand.Intn(len(names))],
			Password:  "password123",
		}

		users[i].Username = strings.ToLower(users[i].FirstName) + strings.ToLower(users[i].LastName) + fmt.Sprintf("%d", i)
		users[i].Email = users[i].Username + "@example.com"
	}
	return users
}

var commentContents = []string{
	"Great post, really enjoyed reading this.",
	"This changed how I think about the topic entirely.",
	"Could you elaborate more on the second point?",
	"I have been looking for something like this for a while.",
	"Interesting perspective, though I slightly disagree.",
	"Well written and easy to follow.",
	"The examples really helped clarify things.",
	"I shared this with my team immediately.",
	"This is exactly what I needed today.",
	"A few of these points are debatable but overall solid.",
	"Would love a follow-up post on this.",
	"Bookmarked. Coming back to this one.",
	"The section on practical application was the best part.",
	"Straightforward and to the point — appreciated.",
	"I learned something new here, thanks.",
	"This resonates a lot with my own experience.",
	"Very well researched.",
	"One of the better articles I have read on this subject.",
	"Not sure I fully agree but it gave me a lot to think about.",
	"The writing style makes a complex topic feel accessible.",
	"Would have liked more real world examples.",
	"Short and sweet, just what I needed.",
	"Really appreciate the depth here.",
	"Sharing this with everyone I know.",
	"This is now my go-to reference on the topic.",
	"Came for the title, stayed for the content.",
	"Thoroughly enjoyed this read.",
	"Could use a bit more nuance in places but still great.",
	"First time commenting — this was worth it.",
	"Keep writing more like this.",
}

func generatePosts(users []*model.User, n int) []*model.Post {
	generatedPosts := make([]*model.Post, n)
	for i := 0; i < n; i++ {
		p := posts[i%len(posts)]
		user := users[rand.Intn(len(users))]
		generatedPosts[i] = &model.Post{
			Title:   p["title"].(string),
			Content: p["content"].(string),
			UserID:  user.ID,
			Tags:    p["tags"].([]string),
		}
	}
	return generatedPosts
}

func generateComments(users []*model.User, posts []*model.Post, n int) []*model.Comment {
	comments := make([]*model.Comment, n)
	for i := 0; i < n; i++ {
		comments[i] = &model.Comment{
			Content: commentContents[rand.Intn(len(commentContents))],
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
		}
	}
	return comments
}
