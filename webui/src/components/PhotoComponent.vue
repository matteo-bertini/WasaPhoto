<script>
	export default {
		props: ["owner", 'photoid', "likesnumber", "commentsnumber", "dateofupload"],
		emits: ["photo_deleted_from_database"],
		// Reactive state
		data() {
			return {
				errormsg: null,
				loading: false,
				PhotoUrl: "",
				PhotoId: "",
				PhotoOwner: "",
				DateOfUpload: "",
				Likes: [],
				LikesNumber: 0,
				Comments: [],
				CommentsNumber: 0,
				Liked: false,
				CommentText: "",
				IsOwner : false
			}
		},
		methods: {

			// Viene eseguito quando viene cliccato il bottone di eliminazione 
			async deletePhoto() {
				let config = {
					headers: {
						"Authorization": `Bearer ${localStorage.getItem("Authstring")}`
					}


				};
				try {
					this.$axios.delete("/users/" + this.owner + "/photos/" + this.photoid + "/", config);

					// Notifico al componente padre che ho eliminato la foto,così la tolga dalla lista 
					this.$emit("photo_deleted_from_database", this.photoid);
					return;

				} catch (e) {
					if (e.response.status == 500) {
						this.$emit("Internal Server Error");
						return;


					}

				}

			},
			async commentPhoto() {
				let config = {
					headers: {
						"Authorization": `Bearer ${localStorage.getItem("Authstring")}`
					},

				};
				let data = { "CommentId": "", "CommentAuthor": `${localStorage.getItem("Username")}`, "CommentText": this.CommentText }
				try {
					let response1 = await this.$axios.post("/users/" + this.PhotoOwner + "/photos/" + this.PhotoId + "/comments/", data, config);
					this.CommentsNumber++;
					this.Comments.unshift(response1.data);

					// Reset del form in cui scrivo il commento 
					this.CommentText = "";

					return;
				} catch (e) {
					//
				}


			},
			async heartButtonClicked() {
				if (this.Liked == false) {
					let config = {
						headers: {
							"Authorization": `Bearer ${localStorage.getItem("Authstring")}`
						},

					};
					try {
						await this.$axios.post("/users/" + this.PhotoOwner + "/photos/" + this.PhotoId + "/likes/", {}, config);
						this.LikesNumber++;
						this.Liked = true;
						this.Likes.unshift(`${localStorage.getItem("Username")}`);

						// Cambio lo stile dell'icona per far vedere che ho messo like
						document.getElementById("hearticon" + this.PhotoId).classList.replace("fa-regular", "fa-solid");
						document.getElementById("heartbutton" + this.PhotoId).style.color = "#8b0000";
						return;
					} catch (e) {
						//
					}
				} else {
					let config = {
						headers: {
							"Authorization": `Bearer ${localStorage.getItem("Authstring")}`
						},

					};
					try {
						await this.$axios.delete("/users/" + this.PhotoOwner + "/photos/" + this.PhotoId + "/likes/" + `${localStorage.getItem("Authstring")}`, config);
						this.LikesNumber--;
						this.Liked = false;
						this.Likes = this.Likes.filter(Like => Like !== `${localStorage.getItem("Username")}`);


						// Cambio lo stile dell'icona per far vedere che ho tolto like
						document.getElementById("hearticon" + this.PhotoId).classList.replace("fa-solid", "fa-regular");
						document.getElementById("heartbutton"+this.PhotoId).style.color = "#000000";
						return;
					} catch (e) {
						//
					}

				}
			},
			updateCommentList(CommentId) {
				this.Comments = this.Comments.filter(Comment => Comment.CommentId !== CommentId);
				this.CommentsNumber--;
				return;

			},
			likePressed(Like) {
				this.$router.replace("/users/"+Like+"/");
				return;
			},
			photoOwnerPressed(){
				this.$router.replace("/users/"+this.PhotoOwner+"/");
				return;

			}
		},


		// Funzione da eseguire appena montato il componente //
		async mounted() {

			// Verrà caricata la foto con tutte le sue informazioni relative aggiornate
			this.PhotoUrl = __API_URL__ + "/users/" + this.owner + "/photos/" + this.photoid + "/";
			this.PhotoOwner = this.owner;
			this.PhotoId = this.photoid;
			this.LikesNumber = this.likesnumber;
			this.CommentsNumber = this.commentsnumber;
			this.DateOfUpload = this.dateofupload;
			this.IsOwner = (`${localStorage.getItem("Username")}`=== this.PhotoOwner );
			try {
				// Ottengo likes e commenti della foto
				let response1 = await this.$axios.get("/users/" + this.PhotoOwner + "/photos/" + this.PhotoId + "/likes/");
				let response2 = await this.$axios.get("/users/" + this.PhotoOwner + "/photos/" + this.PhotoId + "/comments/");
				this.Likes = response1.data.Likes.map(x => x.Username);
				this.Comments = response2.data.Comments;

				// Rendering corretto del pulsante "Mi piace" 
				if (this.Likes.includes(`${localStorage.getItem("Username")}`)) {
					this.Liked = true;
					document.getElementById("hearticon" + this.PhotoId).classList.replace("fa-regular", "fa-solid");
					document.getElementById("heartbutton" + this.PhotoId).style.color = "#8b0000";
					return;
				} else {
					this.Liked = false;
					document.getElementById("hearticon" + this.PhotoId).classList.replace("fa-solid", "fa-regular");
					document.getElementById("heartbutton" + this.PhotoId).style.color = "#000000";
					return;

				}


			} catch (e) {
				if (e.response.status == 500) {
					console.log("Internal Server Error");
					return;
				} else {
					console.log(e);
					return;
				}
			}


		}
	}

</script>

<template>
	<div class="card" style="width: 30em; height: 30em; display: flex; flex-direction: column;">

		<!-- Titolo della card -->
		<div cass="card-title" style="display: flex; flex-direction: row; width: 30em; height: 2em;">
			
			<!-- Owner della foto (solo se chi visualizza non è l'owner della foto) -->
			<div class= "PhotoOwner_div" v-if="!IsOwner" @click="photoOwnerPressed" style="display: flex; flex-direction: row; justify-content: flex-start; align-items: center; margin-top: 0.7em; margin-left:0.9em; width: fit-content; height: 100%;">
				<i class="fa-solid">{{PhotoOwner}}</i>
			</div>

			<!-- Pulsante di eliminazione foto (solo se chi visualizza è l'owner della foto) -->
			<div v-else style="display: flex; flex-direction: row; justify-content: flex-end; width: 100%; height: 100%;">
				<button class="btn"  :disabled = "IsOwner === false" id="deletephotobutton" @click="deletePhoto" type="button" style="border: none;">
					<i class="fa-regular fa-trash-can"></i>
				</button>
			</div>
		</div>

		

		<!-- Foto -->
		<div style="display: flex; flex-direction: row; justify-content: center; height: 70%; width: 100%; margin-top: 15px;">
			<img :src="PhotoUrl" class="card-img-top" alt="Impossibile caricare la foto." data-bs-toggle="modal" :data-bs-target="'#PhotoFullSizeModal'+PhotoId" >
		</div>

		<!-- Pannello Inferiore -->
		<div class="card-body" style="display:flex; flex-direction: column; align-items: center;">
			<div style="display: flex; flex-direction: row; justify-content: center; margin-top: 0px; ">

				<div class="input-group mb-3"
					style=" height: 20px;width: fit-content; display: flex; flex-direction: row; justify-content: center; align-items: center;">

					<!-- Numero di "Mi Piace" -->
					<button class="btn" data-bs-toggle="modal" :data-bs-target="'#LikesModal' +PhotoId"
						style="border:none; background-color: #ffffff; font-size: 20px;">{{LikesNumber}}</button>

					<!-- Pulsante Like -->
					<button class="btn" :id="'heartbutton'+PhotoId" @click="heartButtonClicked" type="button"
						style="border: none;">
						<i class="fa-regular fa-heart" :id="'hearticon'+PhotoId"></i>
					</button>

					<!-- TextArea per inserimento commenti -->
					<input type="text" id="'CommentTextArea'+PhotoId" v-model="CommentText" class="form-control"
						placeholder="Comment..">
					<div class="input-group-append">
						<!-- Pulsante Commenta -->
						<button class="btn" @click="commentPhoto" type="button" :disabled="CommentText.length==0"
							style="border: none;">
							<i class="fa-regular fa-comment"></i>
						</button>
					</div>

					<!-- Numero di Commenti -->
					<button class="btn" data-bs-toggle="modal" :data-bs-target="'#CommentsModal'+PhotoId"
						style="border:none; background-color: white; font-size: 20px;">{{CommentsNumber}}</button>
				</div>
			</div>

		</div>

	</div>

	<!-- LikesModal  -->
	<div class="modal fade" :id="'LikesModal'+PhotoId" tabindex="-1">
		<div class="modal-dialog modal-dialog-scrollable">
			<div class="modal-content">
				<div class="modal-header">
					<h1 class="modal-title fs-5" id="exampleModalLabel">Likes</h1>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body" style="width: fit-content; height: fit-content;">
					<div v-if="this.Likes.length == 0">
						<p>Ancora nessun mi piace</p>
					</div>
					<div v-else>
						<Like data-bs-toggle="modal" :data-bs-target="'#LikesModal'+PhotoId" v-for="Like in Likes" :key ="Like" :LikeUsername="Like" @click="likePressed(Like)">
							{{Like}}
						</Like>
					</div>
					
				</div>
			</div>
		</div>
	</div>

	<!-- CommentsModal  -->
	<div class="modal fade" :id="'CommentsModal' + PhotoId" tabindex="-1">
		<div class="modal-dialog modal-dialog-scrollable">
			<div class="modal-content">
				<div class="modal-header">
					<h1 class="modal-title fs-5" id="exampleModalLabel">Comments</h1>
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body" style="width: fit-content; height: fit-content;">
					<div>
						<div v-if="this.Comments.length == 0">
							<p>Ancora nessun commento</p>

						</div>
						<div v-else>
							<div v-for="Comment in Comments" :key="Comment.CommentId">
								<Comment :PhotoOwner="PhotoOwner" :PhotoId="PhotoId" :CommentId="Comment.CommentId"
									:CommentAuthor="Comment.CommentAuthor" :CommentText="Comment.CommentText"
									@commentdeleted="updateCommentList">
								</Comment>
								<hr />
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- PhotoFullSizeModal -->
	<div class="modal fade" :id="'PhotoFullSizeModal' + PhotoId" tabindex="-1">
		<div class="modal-dialog modal-dialog-scrollable modal-xl" >
			<div class="modal-content" style="display:flex; flex-direction:column; align-items: center;">
				<div class="modal-header" style="display: flex; flex-direction: row; justify-content: flex-end; width: 100%;">
					<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
				</div>
				<div class="modal-body" style="display:flex; flex-direction:row; align-items:center; width: fit-content; height: fit-content;">
					<img :src="PhotoUrl" class="card-img-top" alt="Impossibile caricare la foto.">
				</div>
			</div>
		</div>
	</div>

</template>

<style>
	#deletephotobutton:hover {
		color: #8b0000;
		transform: scale(1.3)
	}
	.btn:hover{
		transform: scale(1.2);
		
	}
	.PhotoOwner_div:hover{
		transform: scale(1.1);
		color:#007bff;
		cursor: pointer;
	}
	
</style>