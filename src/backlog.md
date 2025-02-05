# Backlog
### handin 9 ting vi mangler:
- A transaction must send at least 1 AU to be valid.
- Burde vi lave et variabel, lastBlock, istedet for at altid bare kigge på sidste element i blockslice?
- Checke ordentligt for longest chain

### Generelt:
- ryd op i peer struct
- Genoverveje navne (eg UQTransaction)
- Hvis der er 2 vindere i et slot
- lav lottery om til hjælpe metoder

### Extra hvis vi gider:
- Mange flere test og retter de eksisterende
- Tage højde for hvis en peer modtager en block, som indeholder transactionID'er som peer'en ikke selv har i sin UQTransactions liste.
- Lave main fil, Sebastian mener ikke det er nødvendigt

### Testlist:
- Test that our ten peers have the genesis block
- Test at vores genesis blocks indeholder seeds?
- Test ledgers peers med amount = 1000000
- Test that a block can contain any number of transactions, and might even contain as little as no transactions if there are none to add to the block.

### Kommentarer til rapporten:
- Vi er klar over at peers ikke må kende seeded når de genererer deres public keys, det ville ikke være godt i praksis men det er ligemeget i det her tilfælde
- Vi gemmer alle transactions, også invalid, og smider dem med på træet. De vil aldrig blive udført, fordi hver peer stadig tjekker om de er gyldige eller ej når de skal udføre blocken
- Vi erkender at vi ikke har adresseret følgende problem: i blokken notere vi hvilket slotnummer, peeren vandt i og sender med, men når de andre peers modtager dette er der ikke noget sammenligning for at tjekke om den individuelle peer følger med i samme slot.
- Da vores kode sender alle transactions blokken valid og invalid, og det er op til den enkelte peer at afgøre om en transaction er valid eller ej. Så opstår der en fejl i form af at den individuelle peer i øjeeblikkeet modtager units for hver transaction i den kommende blok. og dette beetyder selvfølgelig at de modtager for mange penge, da de også får fra de invalide transactions i blocken.
