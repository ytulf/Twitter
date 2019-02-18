package main

import (
	"flag"
	"fmt"
	"log"
	//	"time"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func printPageFromFilepath(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if loginPageContent, err := ioutil.ReadFile(filepath); err == nil {
		w.Write(loginPageContent)
	} else {
		log.Fatal("%s not found ! %v", filepath, err)
	}
}

func printFormatedPageFromFilepath(w http.ResponseWriter, filepath string, formatFunc func(string) string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if loginPageContent, err := ioutil.ReadFile(filepath); err == nil {
		w.Write([]byte(formatFunc(string(loginPageContent))))
	} else {
		log.Fatal("%s not found ! %v", filepath, err)
	}
}

func printLoginPage(w http.ResponseWriter) {
	printPageFromFilepath(w, "login_template.html")
}

func setSessionTokenCookie(w *http.ResponseWriter, session_token_value string, req *http.Request) {
	var c http.Cookie
	c.Name = "session_token"
	c.Value = session_token_value
	http.SetCookie(*w, &c)
	req.AddCookie(&c)
}

func redirector(w http.ResponseWriter, req *http.Request, newUrl string) {
	http.Redirect(w, req, newUrl, http.StatusTemporaryRedirect)
}

// callback déclenchée lors d'une requète HTTP vers /
// Son comportement par défaut est une redirection de l'utilisateur vers /login
func loginRedirector(w http.ResponseWriter, req *http.Request) {
	//u := req.URL
	//newUrl := u.Scheme+"://"+u.Host+"/login"
	redirector(w, req, "/login")
}

func getUserNameFromUserid(userId string) string {
	ret := ""
	cmd := exec.Command("/Users/Keijix/Desktop/Twitter/tp", "userName", userId)
	log.Println("getting user name by calling : tp userName", userId)
	out, err := cmd.CombinedOutput()
	if err == nil && out != nil {
		ret += string(out)
	}
	return ret
}

func getTweetsForUserid(userId string) string {
	ret := ""
	cmd := exec.Command("/Users/Keijix/Desktop/Twitter/tp", "listTweetsForUserId", userId)
	log.Println("listing tweets by calling : tp listTweetsForUserId", userId)
	out, err := cmd.CombinedOutput()
	if err == nil && out != nil {
		ret += string(out)
	}
	return ret
}

func recordNewTweet(userId string, tweetContent string) {
	// cmd := fmt.Sprintf("\"/Users/Keijix/Desktop/Twitter/tp recordNewTweet %s %s\"", userId, tweetContent)
	// cmde := exec.Command("/bin/sh", "-c", cmd)
	// log.Println("recording new tweets by calling :", "tp", "recordNewTweet", userId, tweetContent)
	// _, _ = cmde.CombinedOutput()
	log.Println("recording new tweets by calling : tp recordNewTweet", userId, tweetContent)
	cmd := exec.Command("/Users/Keijix/Desktop/Twitter/tp", "recordNewTweet", userId, tweetContent)
	_, _ = cmd.CombinedOutput()
}

func getUserIdFromSession(session_token string) string {
	ret := ""
	cmd := exec.Command("/Users/Keijix/Desktop/Twitter/tp", "userIdFromSession", session_token)
	log.Println("getting user name by calling : tp userIdFromSession", session_token)
	out, err := cmd.CombinedOutput()
	if err == nil && out != nil {
		ret += string(out)
	}
	return ret
}

func listUsers(w http.ResponseWriter, req *http.Request) {
	htmlUserList := ""

	cmd := exec.Command("/Users/Keijix/Desktop/Twitter/tp", "listUserIdComaSeparated")
	log.Println("listing user by calling : tp listUserIdComaSeparated")
	out, err := cmd.CombinedOutput()
	if err == nil && out != nil {
		userids := strings.Split(string(out), ",")
		log.Printf("[listUsers] %v",userids)
		for _, userid := range userids {
			//get username
			htmlUserList += "<br/><a href=/user?id=" + userid + ">"
			htmlUserList += getUserNameFromUserid(userid)
			htmlUserList += "</a><br/>"
		}
	}else{
		log.Printf("[listUsers] trouble %v %v",err, out)
	}

	printFormatedPageFromFilepath(w, "user_list_template.html", func(template string) string { return fmt.Sprintf(template, htmlUserList) })

}

func listUser(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	userId := req.FormValue("id")
	if len(userId) > 0 {
		printFormatedPageFromFilepath(w, "user_template.html", func(template string) string {
			return fmt.Sprintf(template, getUserNameFromUserid(userId), getTweetsForUserid(userId))
		})
	} else {
		w.Write([]byte("Wrong user id"))
	}
}

func login(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	// Si la page de login a déjà été utilisée et
	// que des données sont renvoyées
	if _, ok := req.Form["action"]; ok {
		action_str := req.FormValue("action")
		switch action_str {
		case "register":
			username := req.FormValue("username")
			password := req.FormValue("password")
			passwordCheck := req.FormValue("password_check")
			if password == passwordCheck {
				log.Println("registration using : tp register", username, password)
				log.Println("tp register should return session token for this new user or session token for already registered user")
				cmd := exec.Command("/Users/Keijix/Desktop/Twitter/tp", "register", username, password)
				out, err := cmd.CombinedOutput()
				if err == nil && out != nil && len(string(out)) > 0 {
					setSessionTokenCookie(&w, string(out), req)
					redirector(w, req, "/home")
				} else {
					printLoginPage(w)
				}

			} else {
				w.Write([]byte("Password check Failed"))
				printLoginPage(w)
			}
			break
		case "login":
			username := req.FormValue("username")
			password := req.FormValue("password")
			log.Println("login using : tp login", username, password)
			log.Println("tp login should return session token for this existing user or nothing if there is no user with those credentials")
			cmd := exec.Command("/Users/Keijix/Desktop/Twitter/tp", "login", username, password)
			out, err := cmd.CombinedOutput()
			if err == nil && out != nil && len(string(out)) > 0 {
				setSessionTokenCookie(&w, string(out), req)
				redirector(w, req, "/home")
			} else {
				log.Println("tp Failed to get us correct token")
				// Sinon, affichage de la page de login
				printLoginPage(w)
			}
			break
		}
	} else {
		// Sinon si l'utilisateur est déjà logué redirection vers son home
		if session_token, err := req.Cookie("session_token"); err == nil && len(session_token.Value) > 0 {
			redirector(w, req, "/home")
		} else {
			// Sinon, affichage de la page de login
			printLoginPage(w)
		}
	}
}

func logout(w http.ResponseWriter, req *http.Request) {
	if session_token, err := req.Cookie("session_token"); err == nil {
		log.Println("logout using : tp logout", session_token.Value)
		cmd := exec.Command("/Users/Keijix/Desktop/Twitter/tp", "logout", session_token.Value)
		_, _ = cmd.CombinedOutput()

		session_token.Value = ""
		session_token.MaxAge = 0
		http.SetCookie(w, session_token)
		req.AddCookie(session_token)
	}
	redirector(w, req, "/login")
}

func printUserHomePage(w http.ResponseWriter, req *http.Request) {

	if session_token, err := req.Cookie("session_token"); err == nil && len(session_token.Value) > 1 {
		if userId := getUserIdFromSession(session_token.Value); len(userId) > 0 {
			req.ParseForm()
			// Si la page de login a déjà été utilisée et
			// que des données sont renvoyées
			if _, ok := req.Form["action"]; ok {
				action_str := req.FormValue("action")
				switch action_str {
				case "tweet_compose":
					tweet_content := req.FormValue("tweet_content")
					if len(tweet_content) > 0 {
						recordNewTweet(userId, tweet_content)
					}
					break
				}
			}

			//Affichage du home par défaut
			printFormatedPageFromFilepath(w, "user_home_template.html", func(template string) string {
				// return fmt.Sprintf(template, getUserNameFromUserid(userId), userId)
				// A été inversé, la vrai valeur devrait être :
				return fmt.Sprintf(template, userId, getUserNameFromUserid(userId))
			})
		}
	} else {
		log.Println("[printUserHomePage] You do not have session_token")
		// Sinon, affichage de la page de login
		loginRedirector(w, req)
	}
}

// Défini les différentes routes du frontal web
func defineRoutes() {

	// Pour tout le monde
	http.HandleFunc("/", loginRedirector)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/users", listUsers)
	http.HandleFunc("/user", listUser)

	// Pour les utilisateurs loggués seulement !

	http.HandleFunc("/home", printUserHomePage)
}

// Point d'entrée de l'interface en ligne de commande du serveur web
func main() {
	port := flag.String("p", "192.168.0.12:8100", "host:port to serve on")
	// directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	log.Printf("Serving on HTTP port: %s\n", *port)

	defineRoutes()

	log.Fatal(http.ListenAndServe(*port, nil))
}
