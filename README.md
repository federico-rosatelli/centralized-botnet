---
author:
- Federico Rosatelli 1882771
- Federico Rosatelli 1882771
title:
- 'Sicurezza: Basic Botnet'
- 'Sicurezza: Basic Botnet'
---

Introduzione
============

Si consideri la generica architettura di una botnet composta da agenti
di attacco (bot), installati su dei nodi compromessi, e da un centro di
commando e controllo (C&C). Lo scopo dell'esercitazione è quello di
creare un bot, o meglio detto zombie, in grado di collegarsi al C&C e
ricevere istruzioni sulle azioni da compiere. Da specifica del progetto
questo bot sarà connesso ad un C&C tramite una connessione http in
ascolto su una porta specifica. Nel caso la porta scelta sia già
utilizzata da un altro servizio il bot sarà comunque in grado di
sceglierne una nuova comunicandola al momento della connessione con il
centro di commando e controllo. Lo zombie è stato creato utilizzando
python e oltre 10 librerie diverse per poter eseguire al meglio tutte le
funzionalità richieste. Il C&C invece è un misto di Golang, per la parte
di backend, mongodb come database primario e VueJS per il frontend.

Specifiche e Implementazione
============================

Bot / Zombie
------------

Il bot è un programma compreso da una classe primaria, chiamata Zombie,
e due secondarie per la gestione del server e per il pacchetto di azioni
ricevute dal C&C, chiamate rispettivamente Server e Action. Non appena
avviato il bot inizierà una serie di operazioni in background per
garantire l'utilizzo di tutte le azioni richieste dal centro di
commando. Per prima cosa controlla se è possibile utilizzare comandi
come utente root andando a visionare il gruppo a cui appartiene il
programma e se questo è all'interno dei gruppi sudoers. Come seconda
cosa controlla se è già stato attivato almeno una volta e se si va a
prendere le informazioni della sessione precedente su un file salvato
localmente chiamato config.json. In questo file sono contenute
informazioni come lo stato del programma nell'ultima esecuzione,
l'ultima azione eseguita, le porte in cui era potenzialmente in ascolto
e l'id identificativo dello zombie. Una volta completata la
configurazione, il bot manda le informazioni al centro di commando e
controllo. Le informazioni riguardano lo stato dello zombie, il suo
indirizzo ip, le porte in ascolto e l'identificativo del bot che serve
al C&C per gestire i numerosi e vari zombie sotto il suo controllo. Una
volta mandate queste informazioni, e se il server centrale risponde
autorizzando l'accesso, il bot crea un micro server http in ascolto su
una delle porte specificate accettando solamente richieste di tipo POST
dal C&C. Per far si che solo e solamente il centro di commando e
controllo possa autorizzare le azioni richieste il bot farà una verifica
sulle informazioni in entrata controllando se l'id ricevuto combacia con
il suo. La classe Zombie permette l'implementazione di 3 specifiche
azioni:

1.  Attacco DDOS (Distributed Denial of Services)

2.  Informazioni sul Sistema

3.  Informazioni riguardanti l'uso dell'email tramite il client
    thunderbird

### Attacco DDOS

L'attacco DDOS è implementato dal methodo DDOSAttack() della classe
Zombie. Come argomenti riceve:

-   Il target dell'attacco, ad esempio https://google.com

-   Il tipo di metodo per la richiesta http come GET o POST

-   Il numero di round, quindi quante volte eseguire la richiesta

-   Il payload da aggiungere alla richiesta

-   Il link come fonte di redirect http, ad esempio
    https://example.com/?redirect=your\_target

###  Informazioni sul Sistema

Le informazioni riguardanti tutto il sistema operativo su cui gira il
bot e le caratteristiche hardware sono implementate semplicemente con la
libreria python os. Per poter prendere tutti e soli i dati di interesse
del C&C questa funzionalità è sviluppata dal metodo getSystemInfo() che
prende come argomento una lista di stringhe per la ricerca di
informazioni. Un esempio è la ricerca per version, machine e processor
che in output danno questo ([3](#fig:system_info){reference-type="ref"
reference="fig:system_info"}) tipo di risultato

![[\[fig:system\_info\]]{#fig:system_info
label="fig:system_info"}Esempio di risultato per la ricerca di version,
machine e processor](system_info.png){#fig:system_info
width="0.75\\linewidth"}

### Informazioni sulle Email

Per ricavare le informazioni sulle email personali dell'utente infettato
dal bot il metodo getEmails() ricerca le directory e i file di
configurazione del cliente thunderbird, installato di default in tutti i
sistemi debian/ubuntu. Questi file contengono informazioni molto
preziose come gli indirizzi email utilizzati dall'utente, con le
password a loro associati, la lista dei contatti degli account, le email
inviate e ricevute etc. Queste informazioni, utilizzando un wrapper per
la lettura e scrittura dei file .ini, vengono mandate direttamente al
C&C che le riceve e volendo anche scaricare in locale.

C&C
---

Il centro di commando e controllo è suddiviso, per semplicità in due
parti ben distinte:

-   Backend (Golang)

-   Frontend (VueJS)

La motivazione per l'utilizzo di queste teconologie è la loro
incredibile flessibilità nella gestione e organizzazione dei dati
tramite il protocollo http. Tutte le informazioni sono inserite in un
database NoSql, MongoDb che permette maggiori performance per questi
tipi di dati.

### Backend Golang

Il backend è interamente organizzato e implementato dalla libreria
net/http presente di default come modulo di sistema. Gli entry-points
del server http sono:

-   / in metodo GET

-   /zombie in metodo POST

-   /action in metodo POST

-   /:id/delete in metodo PUT

Il server resta in ascolto sulla porta 8080 accettando tutte le
richieste degli zombie di unirsi alla botnet. Uno zombie che chiede il
suo inserimento sulla botnet manderà una richiesta su /zombie con i dati
relativi al suo indirizzo ip, il suo codice identificativo e la lista
delle porte in ascolto. Se lo zombie non è presente nella rete allora il
C&C lo inserirà nel database creando un nome autogenerato che gli verrà
abbinato. Se l'inserimento è andato a buon fine allora il server manderà
al bot una response positiva comunicandogli di essere stato inserito
correttamente nella botnet In /action il server prende le richieste
dell'utente sul tipo di comando da mandare ai bot. Le richieste devono
essere formattate correttamente in modo da essere eseguite senza errori.
Una tipica richiesta di attacco DDOS ai bot è fornita al server che
riceve, dentro il body, il codice identificativo dello zombie, il target
dell'attaco, il numero di round e altre informazioni secondarie. Il
server C&C che riceve questo tipo di dati andrà a prendere l'indirizzo
ip dello zombie, formatterà i dati, manderà la richiesta al bot e
aspetterà una qualsiasi riposta da esso. Se ci dovesse essere un qualche
tipo di errore il server restituirà quell'errore interrompendo
l'esecuzione del comando. Nel caso il server non sia più in grado di
raggiungere uno zombie, a causa della sua terminazione o altro, il C&C
provvede a settarlo come inattivo temporaneamente. In caso di ripetuti
tentativi di connessione falliti l'utente può cancellare lo zombie
andando a richiamare /:id/zombie, dove :id è il codice identificativo
del bot. Per avere la lista degli zombie presenti nella botnet con il
loro status, l'utente può richiamare / .

### Frontend VueJS

Il Frontend è caratterizzato dalla sola homepage dove l'utente ha la
panoramica di tutti i bot all'interno della botnet. Gli zombie sono
visualizzati in una tabella dove è presente il loro nome, lo status, se
stanno compiendo o meno un'azione, la risposta di tale azione e un tasto
di eliminazione del bot. L'utente può selezionare uno o più zombie alla
volta e comandare un'azione tramite una tendina di azioni possibili
sopra la tabella ([4](#fig:frontend){reference-type="ref"
reference="fig:frontend"})

![[\[fig:frontend\]]{#fig:frontend label="fig:frontend"}Tabella degli
zombie](frontend.png){#fig:frontend width="0.56\\linewidth"}

### Database MongoDb

La scelta del tipo di database è totalmente personale anche se i sistemi
NoSql negli ultimi anni hanno dimostrato una grande robustezza e
capacità nella gestione e archiviazioni di dati, sopratutto di questo
tipo. Il database, chiamato botnet presenta una sola collection,zombie,
dove sono inseriti i dati del bot.

Conclusione
===========

Una botnet di questo tipo è in grado di gestire e comandare in modo
efficiente centinaia di milioni di potenziali zombie sparsi in giro per
il mondo. Gli unici difetti di una rete di questo tipo sono:

1.  Il port forwarding della rete in cui sono presenti gli zombie

2.  Nessun meccanismo di espansione (il bot non è come un virus
    informatico)

3.  La capacità del C&C di rimanere \"nascosto\" e non essere
    rintracciato

Per aggirare questi tipi di problemi, una soluzione protrebbe essere
quella di:

Eliminare il server dallo zombie
--------------------------------

Il bot, a quel punto, può essere visto come un client che aspetta una
risposta continua dal C&C, che altro non sarebbe che l'azione da
compiere, e mettersi in pausa altrimenti.

Aggiungere implementazioni di cloning
-------------------------------------

Lo zombie sarebbe in grado di replicare se stesso andando ad infettare
altre macchine nella stessa rete locale.

Oscuramento dell'indirizzo ip
-----------------------------

Utilizzando vari proxy-server e utilizzando una struttura
decentralizzata si potrebbe \"oscurare\" parzialmente l'indirizzo ip e
le informazioni relative al centro di commando e controllo. L'opzione
migliore sarebbe quella di utlizzare protocolli come tor o i2p.

Introduzione
============

Si consideri la generica architettura di una botnet composta da agenti
di attacco (bot), installati su dei nodi compromessi, e da un centro di
commando e controllo (C&C). Lo scopo dell'esercitazione è quello di
creare un bot, o meglio detto zombie, in grado di collegarsi al C&C e
ricevere istruzioni sulle azioni da compiere. Da specifica del progetto
questo bot sarà connesso ad un C&C tramite una connessione http in
ascolto su una porta specifica. Nel caso la porta scelta sia già
utilizzata da un altro servizio il bot sarà comunque in grado di
sceglierne una nuova comunicandola al momento della connessione con il
centro di commando e controllo. Lo zombie è stato creato utilizzando
python e oltre 10 librerie diverse per poter eseguire al meglio tutte le
funzionalità richieste. Il C&C invece è un misto di Golang, per la parte
di backend, mongodb come database primario e VueJS per il frontend.

Specifiche e Implementazione
============================

Bot / Zombie
------------

Il bot è un programma compreso da una classe primaria, chiamata Zombie,
e due secondarie per la gestione del server e per il pacchetto di azioni
ricevute dal C&C, chiamate rispettivamente Server e Action. Non appena
avviato il bot inizierà una serie di operazioni in background per
garantire l'utilizzo di tutte le azioni richieste dal centro di
commando. Per prima cosa controlla se è possibile utilizzare comandi
come utente root andando a visionare il gruppo a cui appartiene il
programma e se questo è all'interno dei gruppi sudoers. Come seconda
cosa controlla se è già stato attivato almeno una volta e se si va a
prendere le informazioni della sessione precedente su un file salvato
localmente chiamato config.json. In questo file sono contenute
informazioni come lo stato del programma nell'ultima esecuzione,
l'ultima azione eseguita, le porte in cui era potenzialmente in ascolto
e l'id identificativo dello zombie. Una volta completata la
configurazione, il bot manda le informazioni al centro di commando e
controllo. Le informazioni riguardano lo stato dello zombie, il suo
indirizzo ip, le porte in ascolto e l'identificativo del bot che serve
al C&C per gestire i numerosi e vari zombie sotto il suo controllo. Una
volta mandate queste informazioni, e se il server centrale risponde
autorizzando l'accesso, il bot crea un micro server http in ascolto su
una delle porte specificate accettando solamente richieste di tipo POST
dal C&C. Per far si che solo e solamente il centro di commando e
controllo possa autorizzare le azioni richieste il bot farà una verifica
sulle informazioni in entrata controllando se l'id ricevuto combacia con
il suo. La classe Zombie permette l'implementazione di 3 specifiche
azioni:

1.  Attacco DDOS (Distributed Denial of Services)

2.  Informazioni sul Sistema

3.  Informazioni riguardanti l'uso dell'email tramite il client
    thunderbird

### Attacco DDOS

L'attacco DDOS è implementato dal methodo DDOSAttack() della classe
Zombie. Come argomenti riceve:

-   Il target dell'attacco, ad esempio https://google.com

-   Il tipo di metodo per la richiesta http come GET o POST

-   Il numero di round, quindi quante volte eseguire la richiesta

-   Il payload da aggiungere alla richiesta

-   Il link come fonte di redirect http, ad esempio
    https://example.com/?redirect=your\_target

###  Informazioni sul Sistema

Le informazioni riguardanti tutto il sistema operativo su cui gira il
bot e le caratteristiche hardware sono implementate semplicemente con la
libreria python os. Per poter prendere tutti e soli i dati di interesse
del C&C questa funzionalità è sviluppata dal metodo getSystemInfo() che
prende come argomento una lista di stringhe per la ricerca di
informazioni. Un esempio è la ricerca per version, machine e processor
che in output danno questo ([3](#fig:system_info){reference-type="ref"
reference="fig:system_info"}) tipo di risultato

![[\[fig:system\_info\]]{#fig:system_info
label="fig:system_info"}Esempio di risultato per la ricerca di version,
machine e processor](system_info.png){#fig:system_info
width="0.75\\linewidth"}

### Informazioni sulle Email

Per ricavare le informazioni sulle email personali dell'utente infettato
dal bot il metodo getEmails() ricerca le directory e i file di
configurazione del cliente thunderbird, installato di default in tutti i
sistemi debian/ubuntu. Questi file contengono informazioni molto
preziose come gli indirizzi email utilizzati dall'utente, con le
password a loro associati, la lista dei contatti degli account, le email
inviate e ricevute etc. Queste informazioni, utilizzando un wrapper per
la lettura e scrittura dei file .ini, vengono mandate direttamente al
C&C che le riceve e volendo anche scaricare in locale.

C&C
---

Il centro di commando e controllo è suddiviso, per semplicità in due
parti ben distinte:

-   Backend (Golang)

-   Frontend (VueJS)

La motivazione per l'utilizzo di queste teconologie è la loro
incredibile flessibilità nella gestione e organizzazione dei dati
tramite il protocollo http. Tutte le informazioni sono inserite in un
database NoSql, MongoDb che permette maggiori performance per questi
tipi di dati.

### Backend Golang

Il backend è interamente organizzato e implementato dalla libreria
net/http presente di default come modulo di sistema. Gli entry-points
del server http sono:

-   / in metodo GET

-   /zombie in metodo POST

-   /action in metodo POST

-   /:id/delete in metodo PUT

Il server resta in ascolto sulla porta 8080 accettando tutte le
richieste degli zombie di unirsi alla botnet. Uno zombie che chiede il
suo inserimento sulla botnet manderà una richiesta su /zombie con i dati
relativi al suo indirizzo ip, il suo codice identificativo e la lista
delle porte in ascolto. Se lo zombie non è presente nella rete allora il
C&C lo inserirà nel database creando un nome autogenerato che gli verrà
abbinato. Se l'inserimento è andato a buon fine allora il server manderà
al bot una response positiva comunicandogli di essere stato inserito
correttamente nella botnet In /action il server prende le richieste
dell'utente sul tipo di comando da mandare ai bot. Le richieste devono
essere formattate correttamente in modo da essere eseguite senza errori.
Una tipica richiesta di attacco DDOS ai bot è fornita al server che
riceve, dentro il body, il codice identificativo dello zombie, il target
dell'attaco, il numero di round e altre informazioni secondarie. Il
server C&C che riceve questo tipo di dati andrà a prendere l'indirizzo
ip dello zombie, formatterà i dati, manderà la richiesta al bot e
aspetterà una qualsiasi riposta da esso. Se ci dovesse essere un qualche
tipo di errore il server restituirà quell'errore interrompendo
l'esecuzione del comando. Nel caso il server non sia più in grado di
raggiungere uno zombie, a causa della sua terminazione o altro, il C&C
provvede a settarlo come inattivo temporaneamente. In caso di ripetuti
tentativi di connessione falliti l'utente può cancellare lo zombie
andando a richiamare /:id/zombie, dove :id è il codice identificativo
del bot. Per avere la lista degli zombie presenti nella botnet con il
loro status, l'utente può richiamare / .

### Frontend VueJS

Il Frontend è caratterizzato dalla sola homepage dove l'utente ha la
panoramica di tutti i bot all'interno della botnet. Gli zombie sono
visualizzati in una tabella dove è presente il loro nome, lo status, se
stanno compiendo o meno un'azione, la risposta di tale azione e un tasto
di eliminazione del bot. L'utente può selezionare uno o più zombie alla
volta e comandare un'azione tramite una tendina di azioni possibili
sopra la tabella ([4](#fig:frontend){reference-type="ref"
reference="fig:frontend"})

![[\[fig:frontend\]]{#fig:frontend label="fig:frontend"}Tabella degli
zombie](frontend.png){#fig:frontend width="0.56\\linewidth"}

### Database MongoDb

La scelta del tipo di database è totalmente personale anche se i sistemi
NoSql negli ultimi anni hanno dimostrato una grande robustezza e
capacità nella gestione e archiviazioni di dati, sopratutto di questo
tipo. Il database, chiamato botnet presenta una sola collection,zombie,
dove sono inseriti i dati del bot.

Conclusione
===========

Una botnet di questo tipo è in grado di gestire e comandare in modo
efficiente centinaia di milioni di potenziali zombie sparsi in giro per
il mondo. Gli unici difetti di una rete di questo tipo sono:

1.  Il port forwarding della rete in cui sono presenti gli zombie

2.  Nessun meccanismo di espansione (il bot non è come un virus
    informatico)

3.  La capacità del C&C di rimanere \"nascosto\" e non essere
    rintracciato

Per aggirare questi tipi di problemi, una soluzione protrebbe essere
quella di:

Eliminare il server dallo zombie
--------------------------------

Il bot, a quel punto, può essere visto come un client che aspetta una
risposta continua dal C&C, che altro non sarebbe che l'azione da
compiere, e mettersi in pausa altrimenti.

Aggiungere implementazioni di cloning
-------------------------------------

Lo zombie sarebbe in grado di replicare se stesso andando ad infettare
altre macchine nella stessa rete locale.

Oscuramento dell'indirizzo ip
-----------------------------

Utilizzando vari proxy-server e utilizzando una struttura
decentralizzata si potrebbe \"oscurare\" parzialmente l'indirizzo ip e
le informazioni relative al centro di commando e controllo. L'opzione
migliore sarebbe quella di utlizzare protocolli come tor o i2p.
