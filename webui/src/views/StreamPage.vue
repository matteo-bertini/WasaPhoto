<script>
    export default{
        data(){
            return{
                Username:"",
                PhotoStream: []
            }
        },
        async mounted(){

            // Setto l'Username
            this.Username=this.$route.params.Username;

            // Richiesta del photostream da mostrare
            let getMyStream_config = {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem("Authstring")}`
                }
            };
            try{
                let getMyStream_response = await this.$axios.get("/users/"+this.Username+"/",getMyStream_config);
                this.PhotoStream = getMyStream_response.data.PhotoStream;
                return;

            }catch(e){
                console.log(e);
                return;
            }

            

        },
        methods: {
            // Click sul pulsante Profile
			ProfileButtonPressed(){
				this.$router.replace("/users/"+this.Username+"/");
				return;
			},

			// Click sul pulsante Logout
			LogoutButtonPressed(){

				// Pulizia del localstorage
				localStorage.clear();

				// Ritorno alla schermata di Login
				this.$router.replace("/login");
				return;

			},
        }
    }

</script>

<template>
    
    <div id="StreamPageContainer" class="container-fluid" style="display: flex; flex-direction: column; min-width: 100vw; min-height: 100vh;">
        <!-- Pulsanti Profile e Logout-->
        <div style="display: flex; flex-direction: row; justify-content: space-between; margin-left: 3rem; margin-top: 2rem; margin-right: 3rem;">
            
            <!-- Back Button -->
            <div>
                <button class="btn btn-primary" id="ProfileButton" @click="ProfileButtonPressed">
                    <i class="fa solid fa-user"> Profile</i>
                </button>
            </div>
            
            <!-- Logout Button -->
            <div>
                <button class="btn btn-dark" id="LogoutButton" @click="LogoutButtonPressed" style="background-color: darkred;">
                    <i class="fa solid fa-right-from-bracket" style="color: black;"> Logout</i>
                </button>
            </div>
        
        
        </div>

        <!-- Titolo: Username stream -->
        <div style="display: flex; flex-direction: row; justify-content: center; margin-top: 5em;">
            <h1>
                <span class="badge badge-primary" style="background-color: black;">
                    {{Username}}'s stream
                </span>
            </h1>
        </div>

        <!-- PhotoStream -->
        <div v-if="PhotoStream.length>0" style="display: flex; flex-direction: column; align-items: center; gap:2rem; margin-top: 5rem;">
            <Photo v-for="Photo in PhotoStream"
                :key = "Photo.PhotoId"
			    :owner = "Photo.Username"
			    :photoid = "Photo.PhotoId"
			    :likesnumber = "Photo.LikesNumber"
			    :commentsnumber = "Photo.CommentsNumber"
			    :dateofupload = "Photo.DateOfUpload">
            </Photo>

        </div>

        <div v-else style="display: flex; flex-direction: column; align-items: center; justify-content: center; min-height:40em;">
            <i class="fa-solid fa-heart-crack" style="font-size: x-large;"> No content yet! Follow somebody to see photos here.</i>

        </div>

       


    </div>

</template>

<style>
    #StreamPageContainer{
        background: rgb(87,32,122);
		background: -moz-linear-gradient(68deg, rgba(87,32,122,1) 15%, rgba(104,20,138,1) 50%, rgba(87,32,122,1) 85%);
		background: -webkit-linear-gradient(68deg, rgba(87,32,122,1) 15%, rgba(104,20,138,1) 50%, rgba(87,32,122,1) 85%);
		background: linear-gradient(68deg, rgba(87,32,122,1) 15%, rgba(104,20,138,1) 50%, rgba(87,32,122,1) 85%);
		filter: progid:DXImageTransform.Microsoft.gradient(startColorstr="#57207a",endColorstr="#57207a",GradientType=1); 
    }

</style>
